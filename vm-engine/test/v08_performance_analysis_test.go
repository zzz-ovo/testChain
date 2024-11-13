package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"sync"
	"testing"
	"time"

	commonPb "chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/protocol/v2"
	docker_go "chainmaker.org/chainmaker/vm-engine/v2"
)

const (
	performContractVersion = "1.0.0"

	// nolint: varcheck,unused
	timeFormat = "2006-01-02 15:04:05"

	loopNum   = 5
	threadNum = 10
)

var (
	performTxContext protocol.TxSimContext
	mockLogger       protocol.Logger
)

type GroupRunInfo struct {
	groupName string
	startTime int64
	endTime   int64
	taskNum   int64
	runtime   float64 // exec time(uint second)
	tps       float64 // tps   task number / exec time
}

func TestDockerGoPerformance(t *testing.T) {

	//step1: get chainmaker configuration setting from mocked data
	fmt.Printf("=== step 1 load mocked chainmaker configuration file ===\n")
	cmConfig, err := getMockedCMConfig()
	if err != nil {
		log.Fatalf("get the mocked chainmaker configuration failed %v\n", err)
	}

	//step2: generate a docker manager instance
	fmt.Printf("=== step 2 Create docker instance ===\n")
	vmInstanceManager, _ := docker_go.NewInstancesManager(
		chainId,
		newMockHoleLogger(nil, testVMLogName),
		cmConfig,
	)
	mockDockerManager = vmInstanceManager.(*docker_go.InstancesManager)
	mockDockerManager.BlockDurationMgr.AddBlockTxsDuration(blockFingerprint)

	//step3: start docker VM
	fmt.Printf("=== step 3 start Docker VM ===\n")
	dockerContainErr := mockDockerManager.StartVM()
	if dockerContainErr != nil {
		log.Fatalf("start docmer manager instance failed %v\n", dockerContainErr)
	}

	//step4: mock sim context
	fmt.Printf("======step4 mock txContext=======\n")
	performTxContext = InitContextTest()
	mockLogger = newMockHoleLogger(nil, testVMLogName)

	testDeployCut()
	//testCutSave()
	//testDeployZXL()

	performMultipleTxs(loopNum, threadNum)

	//time.Sleep(10 * time.Second)
	fmt.Println("tear down")
	tearDownTest()
}

func performMultipleTxs(loopNum, threadNum int) {
	fmt.Println("--------- Ready to analysis --------------")
	time.Sleep(8 * time.Second)
	fmt.Println("---------- Start -------------------------")

	var totalTps float64

	// declare a slice to store group exec info
	groupResultList := make([]*GroupRunInfo, 0, loopNum)

	for loopIndex := 0; loopIndex < loopNum; loopIndex++ {

		// create a groupRunInfo instance
		groupRunInfo := &GroupRunInfo{
			groupName: "#group" + strconv.Itoa(loopIndex),
			startTime: math.MaxInt64,
			endTime:   math.MinInt64,
			taskNum:   int64(threadNum),
		}

		// create startTimeCh  endTimeCh
		startTimeCh := make(chan int64, threadNum)
		endTimeCh := make(chan int64, threadNum)

		wg := sync.WaitGroup{}

		for threadIndex := 0; threadIndex < threadNum; threadIndex++ {

			wg.Add(1)

			go func(i int) {

				startTimeCh <- time.Now().UnixNano()

				//testCutFindByHash()
				testCutSave()
				//testZXLAddPoint()

				endTimeCh <- time.Now().UnixNano()

				wg.Done()

			}(threadIndex)
		}

		wg.Wait()
		// close all channel
		close(startTimeCh)
		close(endTimeCh)

		// get start time and end time

		// for range startTimeCh (min start time)
		for timeStamp := range startTimeCh {
			if timeStamp < groupRunInfo.startTime {
				groupRunInfo.startTime = timeStamp
			}
		}

		// for range endTimeCh (max end time)
		for timeStamp := range endTimeCh {
			if timeStamp > groupRunInfo.endTime {
				groupRunInfo.endTime = timeStamp
			}
		}

		// compute tps
		groupRunInfo.runtime = float64(groupRunInfo.endTime-groupRunInfo.startTime) / 1e9
		groupRunInfo.tps = cutOutFloat64(float64(groupRunInfo.taskNum)/groupRunInfo.runtime, 2)

		totalTps += groupRunInfo.tps

		fmt.Printf("finished group[%d]: tps:%v\n", loopIndex, groupRunInfo.tps)

		groupResultList = append(groupResultList, groupRunInfo)

	}

	fmt.Println("--------- Finished analysis --------------")
	averageTps := totalTps / float64(loopNum)
	fmt.Println("total tps: ", averageTps)

	//resultFileName := fmt.Sprintf("./result/%s-%d-%d_analysis_result.log", ContractNameTest, loopNum, threadNum)
	//chartFileName := fmt.Sprintf("./result/%s-%d-%d_analysis_chart.html", ContractNameTest, loopNum, threadNum)
	//PrintGroupRunResultAndDraw("result", groupResultList, resultFileName, chartFileName)
	//time.Sleep(3 * time.Second)
}

// ============================= contract_cut functions ===============================

func testDeployCut() {

	fmt.Printf("=== step 6 deploy 【%s】 ===\n", ContractNameTest)

	filePath := fmt.Sprintf("./testdata/%s", ContractNameTest)
	contractBin, contractFileErr := ioutil.ReadFile(filePath)
	if contractFileErr != nil {
		log.Fatal(fmt.Errorf("get byte code failed %v", contractFileErr))
	}

	newRuntimeInstance, _ := mockDockerManager.NewRuntimeInstance(nil, chainId, "",
		"", nil, nil, mockLogger)

	newContractId := &commonPb.Contract{
		Name:        ContractNameTest,
		Version:     performContractVersion,
		RuntimeType: commonPb.RuntimeType_GO,
	}

	parameters := generateInitParams()
	result, _ := newRuntimeInstance.Invoke(newContractId, initMethod, contractBin, parameters,
		performTxContext, uint64(123))
	if result.Code == 0 {
		fmt.Printf("deploy user contract successfully\n")
	}

}

func testCutSave() {
	newRuntimeInstance, _ := mockDockerManager.NewRuntimeInstance(nil, chainId, "",
		"", nil, nil, mockLogger)

	newContractId := &commonPb.Contract{
		Name:        ContractNameTest,
		Version:     performContractVersion,
		RuntimeType: commonPb.RuntimeType_GO,
	}

	parameters := generateInitParams()
	parameters["file_key"] = []byte("key")
	parameters["file_name"] = []byte("name")
	method := "save"

	newRuntimeInstance.Invoke(newContractId, method, nil, parameters,
		performTxContext, uint64(123))
}

// 浮点数切分
func cutOutFloat64(number float64, n int) float64 {
	num := strconv.FormatFloat(number, 'f', 2, 64)
	newNum, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return number
	}
	return newNum
}

//func testCutFindByHash() {
//	newRuntimeInstance, _ := mockDockerManager.NewRuntimeInstance(nil, chainId, "",
//		"", nil, nil, mockLogger)
//
//	newContractId := &commonPb.Contract{
//		Name:        ContractNameTest,
//		Version:     performContractVersion,
//		RuntimeType: commonPb.RuntimeType_GO,
//	}
//
//	parameters := generateInitParams()
//	parameters["method"] = []byte("findByFileHash")
//	parameters["file_key"] = []byte("key")
//
//	newRuntimeInstance.Invoke(newContractId, invokeMethod, nil, parameters,
//		performTxContext, uint64(123))
//}
//
//// =================================== zxl functions =====================================
//
//func testDeployZXL() {
//
//	fmt.Printf("=== step 6 deploy 【%s】 ===\n", ContractNameTest)
//
//	filePath := fmt.Sprintf("./testdata/%s.7z", ContractNameTest)
//	contractBin, contractFileErr := ioutil.ReadFile(filePath)
//	if contractFileErr != nil {
//		log.Fatal(fmt.Errorf("get byte code failed %v", contractFileErr))
//	}
//
//	newRuntimeInstance, _ := mockDockerManager.NewRuntimeInstance(nil, chainId, "",
//		"", nil, nil, mockLogger)
//
//	newContractId := &commonPb.Contract{
//		Name:        ContractNameTest,
//		Version:     performContractVersion,
//		RuntimeType: commonPb.RuntimeType_GO,
//	}
//
//	parameters := generateInitParams()
//	result, _ := newRuntimeInstance.Invoke(newContractId, initMethod, contractBin, parameters,
//		performTxContext, uint64(123))
//	if result.Code == 0 {
//		fmt.Printf("deploy user contract successfully\n")
//	}
//}
//
//func testZXLAddPoint() {
//	newRuntimeInstance, _ := mockDockerManager.NewRuntimeInstance(nil, chainId, "",
//		"", nil, nil, mockLogger)
//
//	newContractId := &commonPb.Contract{
//		Name:        ContractNameTest,
//		Version:     performContractVersion,
//		RuntimeType: commonPb.RuntimeType_GO,
//	}
//
//	parameters := generateInitParams()
//	parameters["method"] = []byte("addPoint")
//	parameters["key"] = []byte("zxl")
//	parameters["value"] = []byte("1")
//
//	newRuntimeInstance.Invoke(newContractId, invokeMethod, nil, parameters,
//		performTxContext, uint64(123))
//}
//

//// print all group run result
//func PrintGroupRunResultAndDraw(title string, groupResultList []*GroupRunInfo, resultFileName, chartFileName string) {
//
//	// if the file not exist, create
//	fileObj, err := os.Create(resultFileName)
//	if err != nil {
//		fmt.Println("generate result file err:", err)
//		return
//	}
//
//	timeList := []string{}
//	tpsValues := []float64{}
//
//	// traverses all executions of the current case
//	for _, groupResult := range groupResultList {
//		// write console
//		//fmt.Printf("\n[group name]: %s\n", groupResult.GroupName)
//		//fmt.Printf("[task umber]: %d\n", groupResult.TaskNum)
//		//fmt.Printf("[start time]: %s\n", time.Unix(groupResult.StartTimeStamp/1e9, 0).Format(timeFormat))
//		//fmt.Printf("[end time]: %s\n", time.Unix(groupResult.EndTimeStamp/1e9, 0).Format(timeFormat))
//		//fmt.Printf("[use time]:%v s\n", groupResult.Runtime)
//		//fmt.Printf("[tps]: %v\n", groupResult.Tps)
//		//fmt.Printf("[err number]: %v\n\n", len(groupResult.Errs))
//
//		// write group info to file
//		fmt.Fprintf(fileObj, "[task number]: %d\n", groupResult.taskNum)
//		fmt.Fprintf(fileObj, "[start time]: %s\n", time.Unix(groupResult.startTime/1e9, 0).Format(timeFormat))
//		fmt.Fprintf(fileObj, "[end time]: %s\n", time.Unix(groupResult.endTime/1e9, 0).Format(timeFormat))
//		fmt.Fprintf(fileObj, "[use time]:%v s\n", groupResult.runtime)
//		fmt.Fprintf(fileObj, "[tps]: %v\n", groupResult.tps)
//		fmt.Fprintf(fileObj, "\n\n")
//
//		// xAxis
//		midTime := groupResult.startTime + ((groupResult.endTime - groupResult.startTime) >> 2)
//		// formate time
//		timeList = append(timeList, time.Unix(int64(midTime/1e9), 0).Format(timeFormat))
//
//		// yAxis
//		tpsValues = append(tpsValues, groupResult.tps)
//
//	}
//
//	// draw picture
//	err = Draw(title, chartFileName, timeList, tpsValues)
//	if err != nil {
//		fmt.Printf("generate picture err: %v", err.Error())
//	}
//}
