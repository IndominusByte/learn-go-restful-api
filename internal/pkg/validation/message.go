package validation

import (
	"fmt"
)

var (
	Default  = "Invalid value."
	Required = "Missing data for required field."
	Email    = "Not a valid email address."
	Min      = "Shorter than minimum length %s."
	Max      = "Longer than maximum length %s."
	Len      = "Length must be equal to %s."
	Gte      = "Must be greater than or equal to %s."
	Gt       = "Must be greater than %s."
	Lte      = "Must be less than or equal to %s."
	Lt       = "Must be less than %s."
	Oneof    = "Must be one of: %s."
)

func SetErrorMessage(ev ErrorValidate) string {
	switch ev.Tag {
	case "required":
		return Required
	case "email":
		return Email
	case "min":
		return fmt.Sprintf(Min, ev.Param)
	case "max":
		return fmt.Sprintf(Max, ev.Param)
	case "len":
		return fmt.Sprintf(Len, ev.Param)
	case "gte":
		return fmt.Sprintf(Gte, ev.Param)
	case "gt":
		return fmt.Sprintf(Gt, ev.Param)
	case "lte":
		return fmt.Sprintf(Lte, ev.Param)
	case "lt":
		return fmt.Sprintf(Lt, ev.Param)
	case "oneof":
		return fmt.Sprintf(Oneof, ev.Param)
	}
	return Default
}
