package example

type StructStruct struct {
	Key   string `default:"key"`
	Value string `default:"value"`
}

type StructStd struct {
	FieldString            string              `default:"test"`
	FieldInt               int                 `default:"1"`
	FieldFloat64           float64             `default:"1.1"`
	FieldBytes             []byte              `default:"test"`
	FieldPBytes            *[]byte             `default:"test"`
	FieldStringSlice       []string            `default:"[test,test2]"`
	FieldStringPSlice      []*string           `default:"[test,test2]"`
	FieldPStringSlice      *[]string           `default:"[test,test2]"`
	FieldPStringPSlice     *[]*string          `default:"[test,test2]"`
	FieldBool              bool                `default:"true"`
	FieldMapStringString   map[string]string   `default:"{key1:value1,key2:value2}"`
	FieldMapPStringPString map[*string]*string `default:"{key1:value1,key2:value2}"`
	FieldMapPBytesPBytes   map[*[]byte]*[]byte `default:"{key1:value1,key2:value2}"`
	FieldMapIntString      map[string]string   `default:"{key1:value1,key2:value2}"`
	FieldMapIntInt         map[int]int         `default:"{1:11,2:22}"`
	FieldMapStringInt      map[string]int      `default:"{value1:11,value2:22}"`
	// FieldStructSlice []StructStruct `default:"[{Key:key,Value:value},{Key:key2,Value:value2}]"`
	// FieldSlice   []int             `default:"[1,2,3]"`
	// FieldStruct StructStruct `default:"{Key:key,Value:value}"`
	// FieldStructSlice []StructStruct `default:"[{Key:key,Value:value}]"`
}
