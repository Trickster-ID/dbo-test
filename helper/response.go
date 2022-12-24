package helper

import "strings"

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

type EmptyObject struct{}

func BuildResponse(data interface{}) Response {
	res := Response{
		Status:  true,
		Message: "OK",
		Error:   nil,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, err string, data interface{}) Response {
	res := Response{
		Status:  false,
		Message: message,
		Error:   strings.Split(err, "\n"),
		Data:    data,
	}
	return res
}
