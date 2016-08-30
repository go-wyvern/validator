package validator

import (
	"net/url"
	"strconv"
	"time"
)

func isSlice(f ValidationFunc, k string, v interface{}, params url.Values, cus bool, args ...interface{}) (bool, error) {
	var is_slice = false
	if vSlice, ok := v.([]interface{}); ok {
		is_slice = true
		for _, subV := range vSlice {
			err := f(k, subV, params, cus, args...)
			if err != nil {
				return is_slice, err
			}
		}
	}
	return is_slice, nil
}

//定义一些规则

func mustLength(k string, v interface{}, params url.Values, cus bool, args ...interface{}) error {
	ok, err := isSlice(mustLength, k, v, params, cus, args...)
	if ok {
		return err
	}
	length := args[0]
	if vString, ok := v.(string); ok {
		if len(vString) != length.(int) {
			pErr := NewParamsError(k, v)
			pErr.Args = args
			return pErr.ErrMustLength(cus)
		}
	}
	return nil
}

func mustMin(k string, v interface{}, params url.Values, cus bool, args ...interface{}) error {
	ok, err := isSlice(mustMin, k, v, params, cus, args...)
	if ok {
		return err
	}
	min := args[0]
	if vInt, ok := v.(int); ok {
		if vInt < min.(int) {
			pErr := NewParamsError(k, v)
			pErr.Args = args
			return pErr.ErrMustMin(cus)
		}
	}
	return nil
}

func mustMax(k string, v interface{}, params url.Values, cus bool, args ...interface{}) error {
	ok, err := isSlice(mustMax, k, v, params, cus, args...)
	if ok {
		return err
	}
	max := args[0]
	if vInt, ok := v.(int); ok {
		if vInt > max.(int) {
			pErr := NewParamsError(k, v)
			pErr.Args = args
			return pErr.ErrMustMax(cus)
		}
	}
	return nil
}

func mustLengthRange(k string, v interface{}, params url.Values, cus bool, args ...interface{}) error {
	ok, err := isSlice(mustLengthRange, k, v, params, cus, args...)
	if ok {
		return err
	}
	min := args[0]
	max := args[1]
	if vString, ok := v.(string); ok {
		if len(vString) < min.(int) || len(vString) > max.(int) {
			pErr := NewParamsError(k, v)
			pErr.Args = args
			return pErr.ErrMustLengthRange(cus)
		}
	}
	return nil
}

func mustValues(k string, v interface{}, params url.Values, cus bool, args ...interface{}) error {
	ok, err := isSlice(mustValues, k, v, params, cus, args...)
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
		pErr := NewParamsError(k, v)
		pErr.Args = args
		return pErr.ErrMustValues(cus)
	}
	return nil
}

func mustTimeLayout(k string, v interface{}, params url.Values, cus bool, args ...interface{}) error {
	ok, err := isSlice(mustTimeLayout, k, v, params, cus, args...)
	if ok {
		return err
	}
	layout := args[0]

	if vString, ok := v.(string); ok {
		_, err := time.Parse(layout.(string), vString)
		if err != nil {
			pErr := NewParamsError(k, v)
			pErr.Args = args
			return pErr.ErrMustTimeLayout(cus)
		}
	}
	return nil
}

func mustLessThan(k string, v interface{}, params url.Values, cus bool, args ...interface{}) error {
	ok, err := isSlice(mustLessThan, k, v, params, cus, args...)
	if ok {
		return err
	}
	field := args[0]
	if vInt, ok := v.(int); ok {
		pInt, err := strconv.Atoi(params.Get(field.(string)))
		if err != nil {
			return NewTextError(err.Error())
		}
		if vInt >= pInt {
			pErr := NewParamsError(k, v)
			pErr.Args = args
			return pErr.ErrMustLessThan(cus)
		}
	}
	return nil
}

func mustLargeThan(k string, v interface{}, params url.Values, cus bool, args ...interface{}) error {
	ok, err := isSlice(mustLargeThan, k, v, params, cus, args...)
	if ok {
		return err
	}
	field := args[0]
	if vInt, ok := v.(int); ok {
		pInt, err := strconv.Atoi(params.Get(field.(string)))
		if err != nil {
			return NewTextError(err.Error())
		}
		if vInt <= pInt {
			pErr := NewParamsError(k, v)
			pErr.Args = args
			return pErr.ErrMustLargeThan(cus)
		}
	}
	return nil
}
