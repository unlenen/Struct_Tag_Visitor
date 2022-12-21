package main
type Root struct {
	Level1  *Level1
	testStr string `foo:"openstack=root_str vcloud=root_str_v"`
	testInt *int   `foo:"openstack=root_int vcloud=root_int_v"`
}

type Level1 struct {
	Level2  *Level2
	TestStr string
	TestInt *int `foo:"openstack=l1_int vcloud=l1_int_v"`
}

type Level2 struct {
	testStr  string `foo:"openstack=l2_str vcloud=l2_str_v"`
	testIntP *int   `foo:"openstack=l2_int_p vcloud=l2_int_p_v"`
	testInt  int    `foo:"lv2_int_default"`
	testBool bool   `foo:"lv2_bool_default"`
}

func main() {

	intLevel1 := 51
	intLevel2 := 61
	intRoot := 41

	root := &Root{
		Level1: &Level1{
			TestStr: "level1-str",
			TestInt: &intLevel1,
			Level2: &Level2{
				testStr:  "level2-str",
				testIntP: &intLevel2,
				testInt:  21,
				testBool: true,
			},
		},
		testStr: "root-str",
		testInt: &intRoot,
	}

	tagExtractor := NewTagExtractor("foo", "openstack")
	properties := tagExtractor.Extract(root)

	fmt.Printf("%v", properties)
}
