package response

const (
	StatusSuccess = "success"
	StatusFailure = "failure"
	StatusError   = "error"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Response {
	return Response{Status: StatusSuccess, Data: data}
}

func BadRequest(data interface{}) Response {
	return Response{Status: StatusFailure, Message: "Bad Request", Data: data}
}

func NotFound() Response {
	return Response{Status: StatusFailure, Message: "Not Found"}
}

func MethodNotAllowed() Response {
	return Response{Status: StatusFailure, Message: "Method Not Allowed"}
}

func InternalServerError() Response {
	return Response{Status: StatusError, Message: "Internal Server Error"}
}
