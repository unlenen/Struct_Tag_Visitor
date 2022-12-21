package main


type TagExtractor struct {
	tagName         string
	defaultProvider string
	properties      *map[string]map[string]string
}

func NewTagExtractor(tagName, defaultProvider string) *TagExtractor {
	properties := make(map[string]map[string]string)
	return &TagExtractor{
		tagName:         tagName,
		defaultProvider: defaultProvider,
		properties:      &properties,
	}
}

func (t *TagExtractor) Extract(dataStruct interface{}) *map[string]map[string]string {
	if dataStruct == nil {
		return t.properties
	}
	dataStructElem := reflect.ValueOf(dataStruct).Elem()
	t.extractFromStruct(dataStructElem)
	return t.properties
}

func (t *TagExtractor) extractFromStruct(dataStruct reflect.Value) {
	for i := 0; i < dataStruct.NumField(); i++ {
		t.extractField(dataStruct, i)
	}
}

func (t *TagExtractor) extractField(dataStruct reflect.Value, index int) {
	typeField := dataStruct.Type().Field(index)
	field := dataStruct.Field(index)
	if typeField.Anonymous {
		t.extractFromStruct(reflect.ValueOf(field.Interface()).Elem())
		return
	}

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return
		}
		field = field.Elem()
		fmt.Printf("Pri Type name:%v , kind:ptr %v\n", typeField.Name, field.Kind())
	} else {
		fmt.Printf("Pri Type name:%v , kind:%v\n", typeField.Name, field.Kind())
	}

	if field.Kind() == reflect.Struct {
		t.extractFromStruct(reflect.ValueOf(field.Interface()))
		return
	}

	if tagValue, ok := typeField.Tag.Lookup(t.tagName); ok {
		if !field.CanSet() { // unaccessable field so hack
			reflectedStruct := reflect.New(dataStruct.Type()).Elem()
			reflectedStruct.Set(dataStruct)
			field = reflectedStruct.Field(index)
			field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		}

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		stringVal := fmt.Sprintf("%v", field.Interface())
		if strings.TrimSpace(stringVal) == "" {
			return
		}

		t.parseTagValue(tagValue, stringVal)
	}

}

func (t *TagExtractor) parseTagValue(tagValue string, fieldValue string) {
	vendors := strings.Split(tagValue, " ") // parse different vendors
	for _, vendor := range vendors {
		data := strings.Split(vendor, "=") // parse vendor name  = property name
		var vendorName string
		var parameterName string
		if len(data) == 1 { //if only contains name take openstack vendor as default
			vendorName = t.defaultProvider
			parameterName = data[0]
		} else {
			vendorName = data[0]
			parameterName = data[1]
		}

		propertiesMap, ok := (*t.properties)[vendorName]
		if !ok { // create map if vendor is not created before
			propertiesMap = make(map[string]string)
			(*t.properties)[vendorName] = propertiesMap
		}
		propertiesMap[parameterName] = fieldValue
	}

}
