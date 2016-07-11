package validator

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
)

func Test_Validate(t *testing.T) {
	Convey("测试参数认证", t, func() {
		params := url.Values{}
		params.Set("test_key","1")

		var defaultValidatorRules ValidatorRules = make(map[string][]Rule)
		r:=new(Rule).SetValidateFunc(MustInt)
		defaultValidatorRules.SetRule("test_key",r)
		err:=Validate(params,defaultValidatorRules)
		So(err,ShouldBeNil())
	})
}
