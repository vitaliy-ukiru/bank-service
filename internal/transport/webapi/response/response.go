package response

import "encoding/json"

type M map[string]any

type OkResponse bool

func (o OkResponse) MarshalJSON() ([]byte, error) {
	var data struct {
		Ok bool
	}

	data.Ok = bool(o)
	return json.Marshal(data)
}

const OkStatus = OkResponse(true)

type Response struct {
	Ok     bool
	Result any
}

type ErrorResponse struct {
	Ok    bool
	Error string
}

func Ok(result any) Response {
	return Response{
		Ok:     true,
		Result: result,
	}
}

func Error(err error) ErrorResponse {
	return ErrorResponse{
		Ok:    false,
		Error: err.Error(),
	}
}

func Fail(err string) ErrorResponse {
	return ErrorResponse{
		Ok:    false,
		Error: err,
	}
}
