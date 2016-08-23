package validator

import "fmt"

type ParamsError struct {
	Key         string
	Value       interface{}
	Args        []interface{}
	Text        string
}

var (
	UnknownParamTpl = "未知的参数:%s"
	RequireParamTpl = "%s是必须的参数"
	MustLengthTpl = "参数[%s]的长度必须为%v"
	MustMinTpl = "参数[%s]的最小值必须大于%v"
	MustMaxTpl = "参数[%s]的最大值必须小于%v"
	MustLengthRangeTpl = "参数[%s]的长度必须为大于%v小于%v"
	MustValuesTpl = "参数[%s]的值必须在%v的范围中"
	MustTimeLayoutTpl = "参数[%s]的格式必须是%v"
	MustLessThanTpl = "参数[%s]的值必须小于参数[%s]"
	MustLargeThanTpl = "参数[%s]的值必须大于参数[%s]"
)
//错误接口
func (p *ParamsError) Error() string {
	return p.Text
}
//位置参数
func (p *ParamsError) ErrUnknownParam() *ParamsError {
	p.Text = fmt.Sprintf(UnknownParamTpl, p.Key)
	return p
}
//用于自定义参数
func (p *ParamsError) CustomErrorText(text string) *ParamsError {
	p.Text = text
	return p
}

func (p *ParamsError) ErrRequireParam() *ParamsError {
	p.Text = fmt.Sprintf(RequireParamTpl, p.Key)
	return p
}

func (p *ParamsError) ErrMustLength() *ParamsError {
	length := p.Args[0]
	p.Text = fmt.Sprintf(MustLengthTpl, p.Key, length)
	return p
}

func (p *ParamsError) ErrMustMin() *ParamsError {
	min := p.Args[0]
	p.Text = fmt.Sprintf(MustMinTpl, p.Key, min)
	return p
}

func (p *ParamsError) ErrMustMax() *ParamsError {
	max := p.Args[0]
	p.Text = fmt.Sprintf(MustMaxTpl, p.Key, max)
	return p
}

func (p *ParamsError) ErrMustLengthRange() *ParamsError {
	min := p.Args[0]
	max := p.Args[1]
	p.Text = fmt.Sprintf(MustLengthRangeTpl, p.Key, min,max)
	return p
}

func (p *ParamsError) ErrMustValues() *ParamsError {
	values := p.Args[0]
	p.Text = fmt.Sprintf(MustValuesTpl, p.Key, values.([]interface{}))
	return p
}

func (p *ParamsError) ErrMustTimeLayout() *ParamsError {
	layout := p.Args[0]
	p.Text = fmt.Sprintf(MustTimeLayoutTpl, p.Key, layout)
	return p
}

func (p *ParamsError) ErrMustLessThan() *ParamsError {
	field := p.Args[0]
	p.Text = fmt.Sprintf(MustLessThanTpl, p.Key, field)
	return p
}

func (p *ParamsError) ErrMustLargeThan() *ParamsError {
	field := p.Args[0]
	p.Text = fmt.Sprintf(MustLargeThanTpl, p.Key, field)
	return p
}

func NewParamsError(k string, v interface{}) *ParamsError {
	pErr := new(ParamsError)
	pErr.Key = k
	pErr.Value = v
	return pErr
}

func NewTextError(text string) *ParamsError {
	pErr := new(ParamsError)
	pErr.Text = text
	return pErr
}

