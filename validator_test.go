package validator

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"fmt"
)

func Test_Validate(t *testing.T) {
	Convey("测试参数认证", t, func() {
		params := url.Values{}
		params.Set("test_key","1")

		defaultValidatorRules:=NewValidator()
		defaultValidatorRules.NewParam("test_key").Require(true).MustInt().MustValues([]interface{}{
			1,
			2,
			3,
		})
		err:=Validator(params,defaultValidatorRules)
		if err!=nil{
			fmt.Println(err.Error())
		}
		//So(err,ShouldBeNil())
	})
}
