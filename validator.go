package validator

import (
	"net/url"
	"reflect"
	"strconv"
	"fmt"
)

type RuleSet interface {
	Require(bool) RuleSet
	MustLength(int, ...error) RuleSet
	MustInt(...error) RuleSet
	MustMin(int, ...error) RuleSet
	MustMax(int, ...error) RuleSet
	MustLengthRange(int, int, ...error) RuleSet
	MustValues([]interface{}, ...error) RuleSet
}

type ValidationFunc  func(string, interface{}, ...interface{}) error

type Validate struct {
	IgnoreUnknownParams bool
	requireParams       []string
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
			Perr := NewParamsError(p, params[p][0])
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
					if valueInterface,ok:=v.valueMap[key];ok{
						err = rule.f(key, valueInterface, rule.args...)
					}else{
						err = rule.f(key, value[0], rule.args...)
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

type ruleSet struct {
	valid     *Validate
	paramName string
	setError  error
}

type rule struct {
	f      ValidationFunc
	args   []interface{}
	errMsg error
}

var _ RuleSet = new(ruleSet)

func (v *Validate) NewParam(paramName string) RuleSet {
	r := new(ruleSet)
	r.paramName = paramName
	r.valid = v
	r.valid.typeMap[paramName] = reflect.String
	r.valid.ruleMap[paramName] = append(r.valid.ruleMap[paramName], *new(rule))
	return r
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
			if value == "true" {
				v.valueMap[key] = true
			} else if value == "false" {
				v.valueMap[key] = false
			} else {
				if Terr, ok := v.typeErrMap[key]; ok {
					return Terr
				}
				return fmt.Errorf("参数[%s]格式错误,参数值必须是bool类型", key)
			}
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
		r.valid.requireParams = append(r.valid.requireParams, r.paramName)
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

func (r *ruleSet) MustValues(values []interface{},errs ...error) RuleSet {
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