package validation

import "fmt"

var (
	Default      = "Invalid value."
	Required     = "Missing data for required field."
	EmailAddress = "Not a valid email address."
	MessageMin   = "Shorter than minimum length %s."
	MessageMax   = "Longer than maximum length %s."
)

func SetErrorMessage(ev ErrorValidate) string {
	switch ev.Tag {
	case "required":
		return Required
	case "email":
		return EmailAddress
	case "min":
		return fmt.Sprintf(MessageMin, ev.Param)
	case "max":
		return fmt.Sprintf(MessageMax, ev.Param)
	}
	return Default
}
