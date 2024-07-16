package response

import "encoding/json"

type M map[string]any

type OkResponse bool

func (o OkResponse) MarshalJSON() ([]byte, error) {
	var data struct {
		Ok bool `json:"ok"`
	}

	data.Ok = bool(o)
	return json.Marshal(data)
}

const OkStatus = OkResponse(true)

type Response struct {
	Ok     bool `json:"ok"`
	Result any  `json:"result"`
}

type ErrorResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
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
