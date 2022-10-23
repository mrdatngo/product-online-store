package util

type ErrorData struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatResponseErrorWithMessage(field string, message string) interface{} {
	return &ErrorData{
		Field:   field,
		Message: message,
	}
}
