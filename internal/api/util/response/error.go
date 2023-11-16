package response

func ErrorResponse(err error) Response {
	return Response{
		Status: false,
		Error:  err.Error(),
	}
}
