package validator

import (
	"fmt"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func LoadCustomBindings() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom tag "strongpwd"
		v.RegisterValidation("strongpwd", strongPassword)
	}
}

func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Example rule: at least 5 chars, one number, one special char
	re := regexp2.MustCompile(`^(?=.*[A-Z])(?=.*[^a-zA-Z0-9]).{5,}$`, 0)
	match, _ := re.MatchString(password)
	return match
}

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	case "strongpwd":
		return "Invalid Password. Password should have at least 5 chars, one number, once special character !"
	}
	return fe.Error()
}
