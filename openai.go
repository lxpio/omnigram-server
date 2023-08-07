package api

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *Response) WithMessage(err string) *Response {
	r.Message = err
	return r
}

type OpenAIModel struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}
