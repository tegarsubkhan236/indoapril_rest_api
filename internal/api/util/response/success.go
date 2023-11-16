package response

func SuccessesResponse(data any, total int64) Response {
	return Response{
		Status: true,
		Data:   data,
		Total:  total,
	}
}

func SuccessResponse(data any) Response {
	return Response{
		Status: true,
		Data:   data,
	}
}
