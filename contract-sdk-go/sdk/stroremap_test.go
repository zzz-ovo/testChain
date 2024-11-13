/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package sdk

/*
func TestNewStoreMap(t *testing.T) {
	m1 := "m1"
	var deep int64 = 1
	storeMap := &StoreMap{
		Name:     m1,
		Depth:    deep,
		HashType: defaultHashType,
	}
	storeMapBytes, _ := json.Marshal(storeMap)
	c := gomock.NewController(t)
	defer c.Finish()
	stubInterface := NewMockSDKInterface(c)
	hashName, _ := storeMap.getHash([]byte(m1), WithHashType(defaultHashType))
	stubInterface.EXPECT().GetStateByte(m1+hashName, "").Return([]byte(""), nil).AnyTimes()
	stubInterface.EXPECT().PutStateByte(m1+hashName, "", storeMapBytes).Return(nil).AnyTimes()

	type args struct {
		name string
		deep int64
		stub SDKInterface
	}
	tests := []struct {
		name    string
		args    args
		want    *StoreMap
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				name: m1,
				deep: deep,
				stub: stubInterface,
			},
			want: &StoreMap{
				Name:     m1,
				Depth:    deep,
				HashType: defaultHashType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStoreMap(tt.args.name, tt.args.deep, tt.args.stub)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStoreMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreMap_Del(t *testing.T) {
	m1 := "m1"
	var deep int64 = 1
	key := "key1"
	c := gomock.NewController(t)
	defer c.Finish()
	sdkInstance := NewMockSDKInterface(c)
	sdkInstance.EXPECT().DelState(m1+key, generateKey(m1, []string{key}))
	type fields struct {
		Name     string
		Deep     int64
		hashType crypto.HashType
	}
	type args struct {
		key  []string
		stub SDKInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Name:     m1,
				Deep:     deep,
				hashType: defaultHashType,
			},
			args: args{
				key:  []string{key},
				stub: sdkInstance,
			},
			wantOk:  true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreMap{
				Name:     tt.fields.Name,
				Depth:    tt.fields.Deep,
				HashType: tt.fields.hashType,
			}
			gotOk, err := c.Del(tt.args.key, tt.args.stub)
			if (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Del() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestStoreMap_Exist(t *testing.T) {
	m1 := "m1"
	var deep int64 = 1
	key := "key1"
	c := gomock.NewController(t)
	defer c.Finish()

	stubInterface := NewMockSDKInterface(c)
	stubInterface.EXPECT().GetStateByte(m1+key, generateKey(m1, []string{key})).Return([]byte(""), nil).AnyTimes()

	type fields struct {
		Name     string
		Deep     int64
		hashType crypto.HashType
	}
	type args struct {
		key  []string
		stub SDKInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Name:     m1,
				Deep:     deep,
				hashType: defaultHashType,
			},
			args: args{
				key:  []string{key},
				stub: stubInterface,
			},
			wantOk:  true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreMap{
				Name:     tt.fields.Name,
				Depth:    tt.fields.Deep,
				HashType: tt.fields.hashType,
			}
			gotOk, err := c.Exist(tt.args.key, tt.args.stub)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Exist() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestStoreMap_Get(t *testing.T) {
	m1 := "m1"
	var deep int64 = 1
	key := "key1"
	value := "value1"
	c := gomock.NewController(t)
	defer c.Finish()

	stubInterface := NewMockSDKInterface(c)
	stubInterface.EXPECT().GetStateByte(m1+key, generateKey(m1, []string{key})).Return([]byte(value), nil).AnyTimes()

	type fields struct {
		Name     string
		Deep     int64
		hashType crypto.HashType
	}
	type args struct {
		key  []string
		stub SDKInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Name:     m1,
				Deep:     deep,
				hashType: defaultHashType,
			},
			args: args{
				key:  []string{key},
				stub: stubInterface,
			},
			want:    []byte(value),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreMap{
				Name:     tt.fields.Name,
				Depth:    tt.fields.Deep,
				HashType: tt.fields.hashType,
			}
			got, err := c.Get(tt.args.key, tt.args.stub)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreMap_Set(t *testing.T) {
	m1 := "m1"
	var deep int64 = 1
	key := "key1"
	value := "value1"
	c := gomock.NewController(t)
	defer c.Finish()

	stubInterface := NewMockSDKInterface(c)
	stubInterface.EXPECT().PutStateByte(m1+key, generateKey(m1, []string{key}), []byte(value)).Return(nil).AnyTimes()

	type fields struct {
		Name     string
		Deep     int64
		hashType crypto.HashType
	}
	type args struct {
		key   []string
		value []byte
		stub  SDKInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Name:     m1,
				Deep:     deep,
				hashType: defaultHashType,
			},
			args: args{
				key:   []string{key},
				value: []byte(value),
				stub:  stubInterface,
			},
			wantOk:  true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreMap{
				Name:     tt.fields.Name,
				Depth:    tt.fields.Deep,
				HashType: tt.fields.hashType,
			}
			gotOk, err := c.Set(tt.args.key, tt.args.value, tt.args.stub)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Set() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestStoreMap_checkMapDeep(t *testing.T) {
	m1 := "m1"
	var deep int64 = 1
	key := "key1"
	type fields struct {
		Name string
		Deep int64
	}
	type args struct {
		key []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Name: m1,
				Deep: deep,
			},
			args:    args{key: []string{key}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreMap{
				Name:  tt.fields.Name,
				Depth: tt.fields.Deep,
			}
			if err := c.checkMapDepth(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("checkMapDeep() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoreMap_generateKey(t *testing.T) {
	m1 := "m1"
	var deep int64 = 1
	key := "key1"
	type fields struct {
		Name     string
		Deep     int64
		hashType crypto.HashType
	}
	type args struct {
		key []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  string
	}{
		{
			name: "t1",
			fields: fields{
				Name:     m1,
				Deep:     deep,
				hashType: defaultHashType,
			},
			args:  args{key: []string{key}},
			want:  m1 + key,
			want1: generateKey(m1, []string{key}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreMap{
				Name:     tt.fields.Name,
				Depth:    tt.fields.Deep,
				HashType: tt.fields.hashType,
			}
			if got, field, err := c.generateKey(tt.args.key); got != tt.want || field != tt.want1 || err != nil {
				t.Errorf("generateKey() = %v, field %v", got, tt.want)
			}
		})
	}
}

func TestStoreMap_save(t *testing.T) {
	m1 := "m1"
	var deep int64 = 1
	c := gomock.NewController(t)
	defer c.Finish()

	storeMap := &StoreMap{
		Name:     m1,
		Depth:    deep,
		HashType: defaultHashType,
	}
	storeMapBytes, _ := json.Marshal(storeMap)
	stubInterface := NewMockSDKInterface(c)
	hashName, _ := storeMap.getHash([]byte(m1))
	stubInterface.EXPECT().PutStateByte(m1+hashName, "", storeMapBytes).Return(nil).AnyTimes()
	type fields struct {
		Name     string
		Deep     int64
		hashType crypto.HashType
	}
	type args struct {
		stub SDKInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Name:     m1,
				Deep:     deep,
				hashType: defaultHashType,
			},
			args: args{stub: stubInterface},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StoreMap{
				Name:     tt.fields.Name,
				Depth:    tt.fields.Deep,
				HashType: tt.fields.hashType,
			}
			if err := c.save(tt.args.stub); (err != nil) != tt.wantErr {
				t.Errorf("save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func generateKey(name string, key []string) string {
	hashName, _ := getHash([]byte(name))
	var storeMapKey = hashName
	for _, k := range key {
		storeMapKey, _ = getHash([]byte(storeMapKey + k))
	}
	return storeMapKey
}

func getHash(data []byte, opt ...Options) (string, error) {
	var hf func() hash.Hash
	ht := new(Hash)
	ht.hashType = defaultHashType
	for _, o := range opt {
		o(ht)
	}
	switch ht.hashType {
	case crypto.HASH_TYPE_SM3:
		hf = sm3.New
	case crypto.HASH_TYPE_SHA256:
		hf = sha256.New
	case crypto.HASH_TYPE_SHA3_256:
		hf = sha3.New256
	default:
		return "", fmt.Errorf("unknown hash algorithm")
	}

	f := hf()

	f.Write(data)
	return hex.EncodeToString(f.Sum(nil)), nil
}
*/
