package utilities

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
)

var (
	v = validator.New()
)

type ErrorResponse struct {
	FailedField string `json:"failedfield"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

type ErrorResponses []*ErrorResponse

func NewErrResponses(r []*ErrorResponse) *ErrorResponses {
	errs := ErrorResponses(r)
	return &errs
}

func init() {
	registerCustomValidation()
}

func registerCustomValidation() {
	// validate.RegisterValidation("SomeTag", func(fl validator.FieldLevel) bool {
	// 	check logic
	// })
}

func (e *ErrorResponses) Error() string {
	var errorString string

	for _, response := range *e {
		if response.FailedField != "" {
			errorString += response.FailedField
		}
		if response.Tag != "" {
			errorString += response.Tag
		}
		if response.Value != "" {
			errorString += response.Value
		}
	}
	return errorString
}

// ValidateStruct는 필드에 제공된 값이 validate 태그에 지정된 조건에 부합하는지 검사합니다.
func ValidateStruct[T any](i T) *ErrorResponses {
	var errors []*ErrorResponse
	err := v.Struct(i)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = fmt.Sprintf("%s", err.Value())
			errors = append(errors, &element)
		}
		return NewErrResponses(errors)
	}
	return nil
}
