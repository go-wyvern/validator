package validator

import "fmt"

type ParamsError struct {
	Key   string
	Value string
	Text  string
}

func (p *ParamsError) Error() string {
	return p.Text
}

func (p *ParamsError) ErrUnknownParam() *ParamsError {
	p.Text = fmt.Sprintf("未知的参数:%s", p.Key)
	return p
}

func (p *ParamsError) CustomErrorText(text string) *ParamsError {
	p.Text = text
	return p
}

func NewParamsError(k, v string) *ParamsError {
	pErr := new(ParamsError)
	pErr.Key = k
	pErr.Value = v
	return pErr
}

