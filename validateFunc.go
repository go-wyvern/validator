package validator

import (
	"strconv"
	"fmt"
)
//定义一些规则
func MustInt(k,v string) error{
	_, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("参数[%s]格式错误,参数值必须是int类型", k)
	}
	return nil
}