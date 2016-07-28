package validator

import (
	"fmt"
	"net/url"
	"time"
)

//定义一些规则

func mustLength(k string, v interface{}, params url.Values, args ...interface{}) error {
	length := args[0]
	if vString, ok := v.(string); ok {
		if len(vString) != length.(int) {
			return fmt.Errorf("参数[%s]的长度必须为%v", k, length)
		}
	} else if vSlice, ok := v.([]interface{}); ok {
		if len(vSlice) != length.(int) {
			return fmt.Errorf("参数[%s]的长度必须为%v", k, length)
		}
	}
	return nil
}

func mustMin(k string, v interface{}, params url.Values, args ...interface{}) error {
	min := args[0]
	if vInt, ok := v.(int); ok {
		if vInt < min.(int) {
			return fmt.Errorf("参数[%s]的最小值必须大于%v", k, min)
		}
	}
	return nil
}

func mustMax(k string, v interface{}, params url.Values, args ...interface{}) error {
	max := args[0]
	if vInt, ok := v.(int); ok {
		if vInt > max.(int) {
			return fmt.Errorf("参数[%s]的最大值必须小于%v", k, max)
		}
	}
	return nil
}

func mustLengthRange(k string, v interface{}, params url.Values, args ...interface{}) error {
	min := args[0]
	max := args[1]
	if vString, ok := v.(string); ok {
		if len(vString) < min.(int) || len(vString) > max.(int) {
			return fmt.Errorf("参数[%s]的长度必须为大于%v小于%v", k, min, max)
		}
	} else if vSlice, ok := v.([]interface{}); ok {
		if len(vSlice) < min.(int) || len(vString) > max.(int) {
			return fmt.Errorf("参数[%s]的长度必须为大于%v小于%v", k, min, max)
		}
	}
	return nil
}

func mustValues(k string, v interface{}, params url.Values, args ...interface{}) error {
	var allNotMatch bool = true
	values := args[0]

	for _, value := range values.([]interface{}) {
		if value == v {
			allNotMatch = false
			break
		}
	}
	if allNotMatch {
		return fmt.Errorf("参数[%s]的长度必须在%v的范围中", k, values.([]interface{}))
	}
	return nil
}

func mustTimeLayout(k string, v interface{}, params url.Values, args ...interface{}) error {
	layout := args[0]

	if vString, ok := v.(string); ok {
		_, err := time.Parse(layout.(string), vString)
		if err != nil {
			return fmt.Errorf("参数[%s]的格式必须是%v", k, layout)
		}
	}
	return nil
}
