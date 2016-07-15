package validator

import "fmt"

type ParamsError struct {
	Key   string
	Value string
	Text  string
}
//错误接口
func (p *ParamsError) Error() string {
	return p.Text
}
//位置参数
func (p *ParamsError) ErrUnknownParam() *ParamsError {
	p.Text = fmt.Sprintf("未知的参数:%s", p.Key)
	return p
}
//用于自定义参数
func (p *ParamsError) CustomErrorText(text string) *ParamsError {
	p.Text = text
	return p
}

func (p *ParamsError) ErrRequireParam() *ParamsError {
	p.Text = fmt.Sprintf("%s是必须的参数", p.Key)
	return p
}

func NewParamsError(k, v string) *ParamsError {
	pErr := new(ParamsError)
	pErr.Key = k
	pErr.Value = v
	return pErr
}

func NewTextError(text string) *ParamsError {
	pErr := new(ParamsError)
	pErr.Text =text
	return pErr
}

