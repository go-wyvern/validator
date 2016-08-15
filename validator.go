package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type RuleSet interface {
	Require(bool) RuleSet
	MustLength(int, ...error) RuleSet
	MustInt(...error) RuleSet
	MustInt64(...error) RuleSet
	MustBool(...error) RuleSet
	MustMin(int, ...error) RuleSet
	MustMax(int, ...error) RuleSet
	MustSeparator(string, reflect.Kind, ...error) RuleSet
	MustLengthRange(int, int, ...error) RuleSet
	MustValues([]interface{}, ...error) RuleSet
	MustTimeLayout(string, ...error) RuleSet
	MustLessThan(string, ...error) RuleSet
	MustLargeThan(string, ...error) RuleSet
	MustFunc(ValidationFunc, []interface{}, ...error) RuleSet
}

const (
	ValidTag = "valid"
	ValidateTag = "validate"
)

type ValidationFunc func(string, interface{}, url.Values, ...interface{}) error

type Validator struct {
	IgnoreUnknownParams bool
	requireParams       []string
	requireUrlParams    []string
	splitChar           string
	ruleMap             map[string][]rule
	valueMap            map[string]interface{}
	defaultValueMap     map[string]interface{}
	typeMap             map[string]reflect.Kind
	elemTypeMap         map[string]reflect.Kind
	typeErrMap          map[string]error
}

func NewValidator() *Validator {
	v := new(Validator)
	v.IgnoreUnknownParams = true
	v.ruleMap = make(map[string][]rule)
	v.typeMap = make(map[string]reflect.Kind)
	v.elemTypeMap = make(map[string]reflect.Kind)
	v.valueMap = make(map[string]interface{})
	v.defaultValueMap = make(map[string]interface{})
	return v
}

func Validate(params url.Values, v *Validator) error {
	v.valueMap = make(map[string]interface{})
	for _, p := range v.requireParams {
		if _, ok := params[p]; !ok {
			Perr := new(ParamsError)
			Perr.Key = p
			Perr.ErrRequireParam()
			return Perr
		}
	}
	for key, value := range params {
		if value[0] = "" {
			continue
		}
		if rules, ok := v.ruleMap[key]; ok {
			err := v.valueCheck(key, value[0])
			if err != nil {
				return err
			}
			for _, rule := range rules {
				if rule.f != nil {
					var err error
					if valueInterface, ok := v.valueMap[key]; ok {
						err = rule.f(key, valueInterface, params, rule.args...)
					} else {
						err = rule.f(key, value[0], params, rule.args...)
					}
					if err != nil {
						if rule.errMsg != nil {
							return rule.errMsg
						} else {
							return err
						}
					}
				}
			}
		} else {
			Perr := NewParamsError(key, value[0])
			Perr.ErrUnknownParam()
			return Perr
		}
	}
	return nil
}

func UrlValidator(params map[string]string, v *Validator) error {
	for _, p := range v.requireUrlParams {
		if _, ok := params[p]; !ok {
			Perr := new(ParamsError)
			Perr.Key = p
			Perr.ErrRequireParam()
			return Perr
		}
	}
	for key, value := range params {
		if value[0] = "" {
			continue
		}
		if rules, ok := v.ruleMap[key]; ok {
			err := v.valueCheck(key, value)
			if err != nil {
				return err
			}
			for _, rule := range rules {
				if rule.f != nil {
					var err error
					if valueInterface, ok := v.valueMap[key]; ok {
						err = rule.f(key, valueInterface, nil, rule.args...)
					} else {
						err = rule.f(key, value, nil, rule.args...)
					}
					if err != nil {
						if rule.errMsg != nil {
							return rule.errMsg
						} else {
							return err
						}
					}
				}
			}
		} else {
			Perr := NewParamsError(key, value)
			Perr.ErrUnknownParam()
			return Perr
		}
	}
	return nil
}

type ruleSet struct {
	valid        *Validator
	is_url_param bool
	paramName    string
	setError     error
}

type rule struct {
	f      ValidationFunc
	args   []interface{}
	errMsg error
}

var _ RuleSet = new(ruleSet)

func (v *Validator) NewParam(paramName string, value ... interface{}) RuleSet {
	r := new(ruleSet)
	r.paramName = paramName
	r.valid = v
	r.valid.typeMap[paramName] = reflect.String
	if len(value) == 1 {
		r.valid.defaultValueMap[paramName] = value[0]
	}
	r.valid.ruleMap[paramName] = append(r.valid.ruleMap[paramName], *new(rule))
	return r
}

func (v *Validator) NewUrlParam(paramName string, value ... interface{}) RuleSet {
	r := new(ruleSet)
	r.paramName = paramName
	r.is_url_param = true
	r.valid = v
	r.valid.typeMap[paramName] = reflect.String
	if len(value) == 1 {
		r.valid.defaultValueMap[paramName] = value[0]
	}
	r.valid.ruleMap[paramName] = append(r.valid.ruleMap[paramName], *new(rule))
	return r
}

func (v *Validator) ValuesToStruct(dst interface{}) error {
	vl := reflect.ValueOf(dst)
	if vl.Kind() != reflect.Ptr || vl.Elem().Kind() != reflect.Struct {
		return NewTextError("interface must be a pointer to struct")
	}
	vl = vl.Elem()
	t := vl.Type()

	for i := 0; i < t.NumField(); i++ {
		if vl.Field(i).Kind() == reflect.Struct {
			st := vl.Field(i).Type()
			sv := vl.Field(i)
			for j := 0; j < st.NumField(); j++ {
				paramName := st.Field(j).Tag.Get(ValidTag)
				if _, ok := v.valueMap[paramName]; ok {
					switch v.typeMap[paramName] {
					case reflect.Int:
						sv.Field(j).SetInt(int64(v.valueMap[paramName].(int)))
					case reflect.Int64:
						sv.Field(j).SetInt(v.valueMap[paramName].(int64))
					case reflect.Bool:
						sv.Field(j).SetBool(v.valueMap[paramName].(bool))
					case reflect.String:
						sv.Field(j).SetString(v.valueMap[paramName].(string))
					case reflect.Slice:
						slicev := reflect.MakeSlice(sv.Field(j).Type(), 0, 0)
						vInterface := v.valueMap[paramName].([]interface{})
						for _, sliceV := range vInterface {
							slicev = reflect.Append(slicev, reflect.ValueOf(sliceV))
						}
						sv.Field(i).Set(slicev)
					}
				} else if _, ok := v.defaultValueMap[paramName]; ok {
					switch v.typeMap[paramName] {
					case reflect.Int:
						sv.Field(j).SetInt(int64(v.defaultValueMap[paramName].(int)))
					case reflect.Int64:
						sv.Field(j).SetInt(v.defaultValueMap[paramName].(int64))
					case reflect.Bool:
						sv.Field(j).SetBool(v.defaultValueMap[paramName].(bool))
					case reflect.String:
						sv.Field(j).SetString(v.defaultValueMap[paramName].(string))
					case reflect.Slice:
						slicev := reflect.MakeSlice(sv.Field(j).Type(), 0, 0)
						vInterface := v.defaultValueMap[paramName].([]interface{})
						for _, sliceV := range vInterface {
							slicev = reflect.Append(slicev, reflect.ValueOf(sliceV))
						}
						sv.Field(i).Set(slicev)
					}
				}
			}
		} else {
			paramName := t.Field(i).Tag.Get(ValidTag)
			if _, ok := v.valueMap[paramName]; ok {
				switch v.typeMap[paramName] {
				case reflect.Int:
					vl.Field(i).SetInt(int64(v.valueMap[paramName].(int)))
				case reflect.Int64:
					vl.Field(i).SetInt(v.valueMap[paramName].(int64))
				case reflect.Bool:
					vl.Field(i).SetBool(v.valueMap[paramName].(bool))
				case reflect.String:
					vl.Field(i).SetString(v.valueMap[paramName].(string))
				case reflect.Slice:
					sv := reflect.MakeSlice(vl.Field(i).Type(), 0, 0)
					vInterface := v.valueMap[paramName].([]interface{})
					for _, sliceV := range vInterface {
						sv = reflect.Append(sv, reflect.ValueOf(sliceV))
					}
					vl.Field(i).Set(sv)
				}
			} else if _, ok := v.defaultValueMap[paramName]; ok {
				switch v.typeMap[paramName] {
				case reflect.Int:
					vl.Field(i).SetInt(int64(v.defaultValueMap[paramName].(int)))
				case reflect.Int64:
					vl.Field(i).SetInt(v.defaultValueMap[paramName].(int64))
				case reflect.Bool:
					vl.Field(i).SetBool(v.defaultValueMap[paramName].(bool))
				case reflect.String:
					vl.Field(i).SetString(v.defaultValueMap[paramName].(string))
				case reflect.Slice:
					sv := reflect.MakeSlice(vl.Field(i).Type(), 0, 0)
					vInterface := v.defaultValueMap[paramName].([]interface{})
					for _, sliceV := range vInterface {
						sv = reflect.Append(sv, reflect.ValueOf(sliceV))
					}
					vl.Field(i).Set(sv)
				}
			}
		}
	}
	return nil
}

func (v *Validator) valueCheck(key, value string) error {
	if pType, ok := v.typeMap[key]; ok {
		var err error
		switch pType {
		case reflect.Int:
			v.valueMap[key], err = strconv.Atoi(value)
			if err != nil {
				if Terr, ok := v.typeErrMap[key]; ok {
					return Terr
				}
				return fmt.Errorf("参数[%s]格式错误,参数值必须是int类型", key)
			}
		case reflect.Int64:
			v.valueMap[key], err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				if Terr, ok := v.typeErrMap[key]; ok {
					return Terr
				}
				return fmt.Errorf("参数[%s]格式错误,参数值必须是int64类型", key)
			}
		case reflect.Bool:
			if valueBool, err := strconv.ParseBool(value); err != nil {
				if Terr, ok := v.typeErrMap[key]; ok {
					return Terr
				}
				return fmt.Errorf("参数[%s]格式错误,参数值必须是bool类型", key)
			} else {
				v.valueMap[key] = valueBool
			}
		case reflect.String:
			v.valueMap[key] = value
		case reflect.Slice:
			var sliceInterface []interface{}
			sliceString := strings.Split(value, v.splitChar)
			switch v.elemTypeMap[key] {
			case reflect.Int:
				for _, vString := range sliceString {
					vInt, err := strconv.Atoi(vString)
					if err != nil {
						if Terr, ok := v.typeErrMap[key]; ok {
							return Terr
						}
						return fmt.Errorf("参数[%s]格式错误,参数值必须是int类型", key)
					}
					sliceInterface = append(sliceInterface, vInt)
				}
			case reflect.Int64:
				for _, vString := range sliceString {
					vInt64, err := strconv.ParseInt(vString, 10, 64)
					if err != nil {
						if Terr, ok := v.typeErrMap[key]; ok {
							return Terr
						}
						return fmt.Errorf("参数[%s]格式错误,参数值必须是int64类型", key)
					}
					sliceInterface = append(sliceInterface, vInt64)
				}
			case reflect.Bool:
				for _, vString := range sliceString {
					if vBool, err := strconv.ParseBool(vString); err != nil {
						if Terr, ok := v.typeErrMap[key]; ok {
							return Terr
						}
						return fmt.Errorf("参数[%s]格式错误,参数值必须是bool类型", key)
					} else {
						sliceInterface = append(sliceInterface, vBool)
					}
				}
			case reflect.String:
				for _, vString := range sliceString {
					sliceInterface = append(sliceInterface, vString)
				}
			}
			v.valueMap[key] = sliceInterface
		default:

		}

	}
	return nil
}

func (r *ruleSet) Require(require bool) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set Require")
	}

	if require {
		if r.is_url_param {
			r.valid.requireUrlParams = append(r.valid.requireUrlParams, r.paramName)
		} else {
			r.valid.requireParams = append(r.valid.requireParams, r.paramName)
		}
	}
	return r
}

func (r *ruleSet) MustInt(errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustInt")
	}

	if len(errs) == 1 {
		r.valid.typeErrMap[r.paramName] = errs[0]
	}
	r.valid.typeMap[r.paramName] = reflect.Int
	return r
}

func (r *ruleSet) MustInt64(errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustInt")
	}

	if len(errs) == 1 {
		r.valid.typeErrMap[r.paramName] = errs[0]
	}
	r.valid.typeMap[r.paramName] = reflect.Int64
	return r
}

func (r *ruleSet) MustBool(errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustBool")
	}

	if len(errs) == 1 {
		r.valid.typeErrMap[r.paramName] = errs[0]
	}
	r.valid.typeMap[r.paramName] = reflect.Bool
	return r
}

func (r *ruleSet) MustSeparator(s string, elemType reflect.Kind, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustBool")
	}

	if len(errs) == 1 {
		r.valid.typeErrMap[r.paramName] = errs[0]
	}
	r.valid.splitChar = s
	r.valid.typeMap[r.paramName] = reflect.Slice
	r.valid.elemTypeMap[r.paramName] = elemType
	return r
}

func (r *ruleSet) MustLength(length int, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustLength")
	}
	rl := new(rule)
	rl.f = mustLength
	rl.args = append(rl.args, length)
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)
	return r
}

func (r *ruleSet) MustMin(min int, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustMin")
	}
	rl := new(rule)
	rl.f = mustMin
	rl.args = append(rl.args, min)
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)
	return r
}

func (r *ruleSet) MustMax(max int, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustMax")
	}
	rl := new(rule)
	rl.f = mustMax
	rl.args = append(rl.args, max)
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)
	return r
}

func (r *ruleSet) MustLengthRange(min, max int, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustLengthRange")
	}
	rl := new(rule)
	rl.f = mustLengthRange
	rl.args = append(rl.args, min)
	rl.args = append(rl.args, max)
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)
	return r
}

func (r *ruleSet) MustValues(values []interface{}, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustLengthRange")
	}
	rl := new(rule)
	rl.f = mustValues
	rl.args = append(rl.args, values)
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)
	return r
}

func (r *ruleSet) MustTimeLayout(layout string, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustLengthRange")
	}
	rl := new(rule)
	rl.f = mustTimeLayout
	rl.args = append(rl.args, layout)
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)
	return r
}

func (r *ruleSet) MustLessThan(field string, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustLessThan")
	}
	rl := new(rule)
	rl.f = mustLessThan
	rl.args = append(rl.args, field)
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)
	return r
}

func (r *ruleSet) MustLargeThan(field string, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustLessThan")
	}
	rl := new(rule)
	rl.f = mustLargeThan
	rl.args = append(rl.args, field)
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)
	return r
}

func (r *ruleSet) MustFunc(f ValidationFunc, args []interface{}, errs ...error) RuleSet {
	if r.setError != nil {
		return r
	}
	if r.paramName == "" {
		panic("unknown param name when set MustLengthRange")
	}
	rl := new(rule)
	rl.f = f
	rl.args = args
	if len(errs) == 1 {
		rl.errMsg = errs[0]
	}
	r.valid.ruleMap[r.paramName] = append(r.valid.ruleMap[r.paramName], *rl)

	return r
}
