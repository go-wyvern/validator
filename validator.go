package validator

import "net/url"

var IgnoreUnknownParam bool

type ValidatorRules map[string][]Rule

var defaultValidatorRules ValidatorRules = make(map[string][]Rule)

type ValidationFunc func(k, v string) error

type Rule struct {
	f          ValidationFunc
	errMessage string
}

func Validate(params url.Values, validators ValidatorRules) error {
	for k, v := range params {
		if rules, ok := validators[k]; ok {
			for _, rule := range rules {
				err := rule.f(k, v[0])
				if err != nil {
					if rule.errMessage != "" {
						pErr := NewParamsError(k, v[0])
						return pErr.CustomErrorText(rule.errMessage)
					} else {
						return err
					}
				}
			}
		} else {
			if !IgnoreUnknownParam {
				pErr := NewParamsError(k, v[0])
				return pErr.ErrUnknownParam()
			}
		}
	}
	return nil
}

func (r *Rule) SetValidateFunc(f ValidationFunc) *Rule {
	r.f = f
	return r
}

func (r *Rule) Error(errStr string) *Rule {
	r.errMessage = errStr
	return r
}

func SetRule(paramName string, r *Rule) {
	defaultValidatorRules.SetRule(paramName, r)
}

func (c ValidatorRules) SetRule(paramName string, r *Rule) {
	c[paramName] = append(c[paramName], *r)
}

