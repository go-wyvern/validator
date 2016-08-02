package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type RuleSet interface {
	Require(bool) RuleSet
	MustLength(int, ...error) RuleSet
	MustInt(...error) RuleSet
	MustBool(...error) RuleSet
	MustMin(int, ...error) RuleSet
	MustMax(int, ...error) RuleSet
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

type Validate struct {
	IgnoreUnknownParams bool
	requireParams       []string
	requireUrlParams    []string
	ruleMap             map[string][]rule
	valueMap            map[string]interface{}
	typeMap             map[string]reflect.Kind
	typeErrMap          map[string]error
}

func NewValidator() *Validate {
	v := new(Validate)
	v.IgnoreUnknownParams = true
	v.ruleMap = make(map[string][]rule)
	v.typeMap = make(map[string]reflect.Kind)
	v.valueMap = make(map[string]interface{})
	return v
}

func Validator(params url.Values, v *Validate) error {
	for _, p := range v.requireParams {
		if _, ok := params[p]; !ok {
			Perr := new(ParamsError)
			Perr.Key = p
			Perr.ErrRequireParam()
			return Perr
		}
	}
	for key, value := range params {
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

func UrlValidator(params map[string]string, v *Validate) error {
	for _, p := range v.requireUrlParams {
		if _, ok := params[p]; !ok {
			Perr := new(ParamsError)
			Perr.Key = p
			Perr.ErrRequireParam()
			return Perr
		}
	}
	for key, value := range params {
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
	valid        *Validate
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

func (v *Validate) NewParam(paramName string, value ... interface{}) RuleSet {
	r := new(ruleSet)
	r.paramName = paramName
	r.valid = v
	r.valid.typeMap[paramName] = reflect.String
	if len(value) == 1 {
		r.valid.valueMap[paramName] = value[0]
	}
	r.valid.ruleMap[paramName] = append(r.valid.ruleMap[paramName], *new(rule))
	return r
}

func (v *Validate) NewUrlParam(paramName string, value ... interface{}) RuleSet {
	r := new(ruleSet)
	r.paramName = paramName
	r.is_url_param = true
	r.valid = v
	r.valid.typeMap[paramName] = reflect.String
	if len(value) == 1 {
		r.valid.valueMap[paramName] = value[0]
	}
	r.valid.ruleMap[paramName] = append(r.valid.ruleMap[paramName], *new(rule))
	return r
}

func (v *Validate) ValuesToStruct(dst interface{}) error {
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
					case reflect.Bool:
						sv.Field(j).SetBool(v.valueMap[paramName].(bool))
					case reflect.String:
						sv.Field(j).SetString(v.valueMap[paramName].(string))
					}
				}
			}
		} else {
			paramName := t.Field(i).Tag.Get(ValidTag)
			if _, ok := v.valueMap[paramName]; ok {
				switch v.typeMap[paramName] {
				case reflect.Int:
					vl.Field(i).SetInt(int64(v.valueMap[paramName].(int)))
				case reflect.Bool:
					vl.Field(i).SetBool(v.valueMap[paramName].(bool))
				case reflect.String:
					vl.Field(i).SetString(v.valueMap[paramName].(string))
				}
			}
		}
	}
	return nil
}

func (v *Validate) valueCheck(key, value string) error {
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

func (r *ruleSet) MustLessThan(field string,errs ...error) RuleSet {
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

func (r *ruleSet) MustLargeThan(field string,errs ...error) RuleSet {
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
