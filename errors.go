package validator

import "bytes"

type ParamsError struct {
	Key   string
	Value interface{}
	Args  []interface{}
	Text  string
}

var (
	defaultUnknownParamTpl    = "未知的参数:{{.Key}}"
	defaultRequireParamTpl    = "{{.Key}}是必须的参数"
	defaultRequireNotNullTpl  = "{{.Key}}是必须的参数，不能为空"
	defaultMustLengthTpl      = "参数[{{.Key}}]的长度必须为{{index .Args 0}}"
	defaultMustMinTpl         = "参数[{{.Key}}]的最小值必须大于{{index .Args 0}}"
	defaultMustMaxTpl         = "参数[{{.Key}}]的最大值必须小于{{index .Args 0}}"
	defaultMustLengthRangeTpl = "参数[{{.Key}}]的长度必须为大于{{index .Args 0}}小于{{index .Args 1}}"
	defaultMustValuesTpl      = `参数[{{.Key}}]的值必须在{{index .Args 0}}的范围中`
	defaultMustTimeLayoutTpl  = "参数[{{.Key}}]的格式必须是{{index .Args 0}}"
	defaultMustLessThanTpl    = "参数[{{.Key}}]的值必须小于参数[{{index .Args 0}}]"
	defaultMustLargeThanTpl   = "参数[{{.Key}}]的值必须大于参数[{{index .Args 0}}]"
)

var (
	CustomUnknownParamTpl    = "{{.unknow_param}}"
	CustomRequireParamTpl    = "{{.require_param}}"
	CustomRequireNotNullTpl  = "{{.require_not_null}}"
	CustomMustLengthTpl      = "{{.must_length}}"
	CustomMustMinTpl         = "{{.must_min}}"
	CustomMustMaxTpl         = "{{.must_max}}"
	CustomMustLengthRangeTpl = "{{.must_length_range}}"
	CustomMustValuesTpl      = "{{.must_values}}"
	CustomMustTimeLayoutTpl  = "{{.must_time_layout}}"
	CustomMustLessThanTpl    = "{{.must_less_than}}"
	CustomMustLargeThanTpl   = "{{.must_large_than}}"
)

//错误接口
func (p *ParamsError) Error() string {
	return p.Text
}

//位置参数
func (p *ParamsError) ErrUnknownParam(cus bool) *ParamsError {
	if cus {
		p.Text = CustomUnknownParamTpl
		return p
	}
	p.Text = defaultUnknownParamTpl
	return p.Tr()
}

//用于自定义参数
func (p *ParamsError) CustomErrorText(text string) *ParamsError {
	p.Text = text
	return p
}

func (p *ParamsError) ErrRequireParam(cus bool) *ParamsError {
	if cus {
		p.Text = CustomRequireParamTpl
		return p
	}
	p.Text = defaultRequireParamTpl
	return p.Tr()
}

func (p *ParamsError) ErrRequireNotNull(cus bool) *ParamsError {
	if cus {
		p.Text = CustomRequireNotNullTpl
		return p
	}
	p.Text = defaultRequireNotNullTpl
	return p.Tr()
}

func (p *ParamsError) ErrMustLength(cus bool) *ParamsError {
	if cus {
		p.Text = CustomMustLengthTpl
		return p
	}
	p.Text = defaultMustLengthTpl
	return p.Tr()
}

func (p *ParamsError) ErrMustMin(cus bool) *ParamsError {
	if cus {
		p.Text = CustomMustMinTpl
		return p
	}
	p.Text = defaultMustMinTpl
	return p.Tr()
}

func (p *ParamsError) ErrMustMax(cus bool) *ParamsError {
	if cus {
		p.Text = CustomMustMaxTpl
		return p
	}
	p.Text = defaultMustMaxTpl
	return p.Tr()
}

func (p *ParamsError) ErrMustLengthRange(cus bool) *ParamsError {
	if cus {
		p.Text = CustomMustLengthRangeTpl
		return p
	}
	p.Text = defaultMustLengthRangeTpl
	return p.Tr()
}

func (p *ParamsError) ErrMustValues(cus bool) *ParamsError {
	if cus {
		p.Text = CustomMustValuesTpl
		return p
	}
	p.Text = defaultMustValuesTpl
	return p.Tr()
}

func (p *ParamsError) ErrMustTimeLayout(cus bool) *ParamsError {
	if cus {
		p.Text = CustomMustTimeLayoutTpl
		return p
	}
	p.Text = defaultMustTimeLayoutTpl
	return p.Tr()
}

func (p *ParamsError) ErrMustLessThan(cus bool) *ParamsError {
	if cus {
		p.Text = CustomMustLessThanTpl
		return p
	}
	p.Text = defaultMustLessThanTpl
	return p.Tr()
}

func (p *ParamsError) ErrMustLargeThan(cus bool) *ParamsError {
	if cus {
		p.Text = CustomMustLargeThanTpl
		return p
	}
	p.Text = defaultMustLargeThanTpl
	return p.Tr()
}

func (p *ParamsError) Tr() *ParamsError {
	var buff = make([]byte, 0)
	b := bytes.NewBuffer(buff)
	tmpl(b, p.Text, p)
	p.Text = b.String()
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
