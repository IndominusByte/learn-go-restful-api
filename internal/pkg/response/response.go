package response

import (
	"encoding/json"
	"net/http"
)

var responses = map[int]map[string]interface{}{
	200: {
		"status_code": 200,
		"status":      true,
		"message":     "Request Success.",
		"results":     nil,
	},
	201: {
		"status_code": 201,
		"status":      true,
		"message":     "Request Created.",
	},
	400: {
		"status_code":    400,
		"status":         false,
		"message":        "Validation Failed.",
		"detail_message": nil,
		"results":        nil,
	},
	409: {
		"status_code":    409,
		"status":         false,
		"message":        "Conflict.",
		"detail_message": nil,
		"results":        nil,
	},
	401: {
		"status_code":    401,
		"status":         false,
		"message":        "Unauthorized Request.",
		"detail_message": nil,
		"results":        nil,
	},
	403: {
		"status_code":    403,
		"status":         false,
		"message":        "Forbidden Request.",
		"detail_message": nil,
		"results":        nil,
	},
	404: {
		"status_code":    404,
		"status":         false,
		"message":        "Resource Not Found.",
		"detail_message": nil,
		"results":        nil,
	},
	413: {
		"status_code":    413,
		"status":         false,
		"message":        "Payload Too Large.",
		"detail_message": nil,
		"results":        nil,
	},
	422: {
		"status_code":    422,
		"status":         false,
		"message":        "Unprocessable Entity.",
		"detail_message": nil,
		"results":        nil,
	},
	500: {
		"status_code":    500,
		"status":         false,
		"message":        "Internal Server Error.",
		"detail_message": nil,
		"results":        nil,
	},
}

func WriteJSONResponse(rw http.ResponseWriter, code int, results interface{}, detail_message interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(code)
	if response := responses[code]; response != nil {
		response["results"] = results
		response["detail_message"] = detail_message
		json.NewEncoder(rw).Encode(response)
		return
	}
	panic("Response code not found!")
}
