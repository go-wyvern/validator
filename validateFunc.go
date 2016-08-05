package validator

import (
	"fmt"
	"net/url"
	"time"
	"strconv"
)

func isSlice(f ValidationFunc, k string, v interface{}, params url.Values, args ...interface{}) (bool, error) {
	var is_slice = false
	if vSlice, ok := v.([]interface{}); ok {
		is_slice = true
		for _, subV := range vSlice {
			err := f(k, subV, params, args...)
			if err != nil {
				return is_slice, err
			}
		}
	}
	return is_slice, nil
}

//定义一些规则

func mustLength(k string, v interface{}, params url.Values, args ...interface{}) error {
	ok, err := isSlice(mustLength, k, v, params, args...)
	if ok {
		return err
	}
	length := args[0]
	if vString, ok := v.(string); ok {
		if len(vString) != length.(int) {
			return fmt.Errorf("参数[%s]的长度必须为%v", k, length)
		}
	}
	return nil
}

func mustMin(k string, v interface{}, params url.Values, args ...interface{}) error {
	ok, err := isSlice(mustMin, k, v, params, args...)
	if ok {
		return err
	}
	min := args[0]
	if vInt, ok := v.(int); ok {
		if vInt < min.(int) {
			return fmt.Errorf("参数[%s]的最小值必须大于%v", k, min)
		}
	}
	return nil
}

func mustMax(k string, v interface{}, params url.Values, args ...interface{}) error {
	ok, err := isSlice(mustMax, k, v, params, args...)
	if ok {
		return err
	}
	max := args[0]
	if vInt, ok := v.(int); ok {
		if vInt > max.(int) {
			return fmt.Errorf("参数[%s]的最大值必须小于%v", k, max)
		}
	}
	return nil
}

func mustLengthRange(k string, v interface{}, params url.Values, args ...interface{}) error {
	ok, err := isSlice(mustLengthRange, k, v, params, args...)
	if ok {
		return err
	}
	min := args[0]
	max := args[1]
	if vString, ok := v.(string); ok {
		if len(vString) < min.(int) || len(vString) > max.(int) {
			return fmt.Errorf("参数[%s]的长度必须为大于%v小于%v", k, min, max)
		}
	}
	return nil
}

func mustValues(k string, v interface{}, params url.Values, args ...interface{}) error {
	ok, err := isSlice(mustValues, k, v, params, args...)
	if ok {
		return err
	}
	var allNotMatch bool = true
	values := args[0]

	for _, value := range values.([]interface{}) {
		if value == v {
			allNotMatch = false
			break
		}
	}

	if allNotMatch {
		return fmt.Errorf("参数[%s]的值必须在%v的范围中", k, values.([]interface{}))
	}
	return nil
}

func mustTimeLayout(k string, v interface{}, params url.Values, args ...interface{}) error {
	ok, err := isSlice(mustTimeLayout, k, v, params, args...)
	if ok {
		return err
	}
	layout := args[0]

	if vString, ok := v.(string); ok {
		_, err := time.Parse(layout.(string), vString)
		if err != nil {
			return fmt.Errorf("参数[%s]的格式必须是%v", k, layout)
		}
	}
	return nil
}

func mustLessThan(k string, v interface{}, params url.Values, args ...interface{}) error {
	ok, err := isSlice(mustLessThan, k, v, params, args...)
	if ok {
		return err
	}
	field := args[0]
	if vInt, ok := v.(int); ok {
		pInt, err := strconv.Atoi(params.Get(field.(string)))
		if err != nil {
			return err
		}
		if vInt >= pInt {
			return fmt.Errorf("参数[%s]的值必须小于参数[%s]", k, field)
		}
	}
	return nil
}

func mustLargeThan(k string, v interface{}, params url.Values, args ...interface{}) error {
	ok, err := isSlice(mustLargeThan, k, v, params, args...)
	if ok {
		return err
	}
	field := args[0]
	if vInt, ok := v.(int); ok {
		pInt, err := strconv.Atoi(params.Get(field.(string)))
		if err != nil {
			return err
		}
		if vInt <= pInt {
			return fmt.Errorf("参数[%s]的值必须大于参数[%s]", k, field)
		}
	}
	return nil
}
