## Package validator

Package validator implements variable validations

## Installation

Just use go get.

```
go get github.com/go-wyvern/validator
```

And then just import the package into your own code.

```
import (
    "github.com/go-wyvern/validator"
)
```

## Usage

```
params := url.Values{}
params.Set("foo","200")
params.Set("bar","active")

defaultValidator:=NewValidator()
defaultValidator.NewParam("foo").Require(true).MustInt().MustMax(256).MustMin(128)
defaultValidator.NewParam("bar").Require(true).MustInt().MustMax(10).MustMin(3).MustValues([]interface{}{
	"active",
	"unactive",
})		
err:=Validator(params,defaultValidatorRules)
```

Here is the list of validators interface in the package. 

```
type RuleSet interface {
	Require(bool) RuleSet
	MustLength(int, ...error) RuleSet
	MustInt(...error) RuleSet
	MustMin(int, ...error) RuleSet
	MustMax(int, ...error) RuleSet
	MustLengthRange(int, int, ...error) RuleSet
	MustValues([]interface{}, ...error) RuleSet
}
```

## License
Copyright 2016 Wyvern wyvern.wu@aliyun.com
