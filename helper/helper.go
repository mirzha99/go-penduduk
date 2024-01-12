package helper

// SuccessResponse adalah struktur respons sukses
type SuccessResponse struct {
	StatusCode int    `json:"statuscode"`
	Result     any    `json:"result"`
	Message    string `json:"message"`
	Token      string `json:"token"`
}

func (s *SuccessResponse) SuccesRMessage() any {
	return map[string]any{
		"statuscode": s.StatusCode,
		"message":    s.Message,
		"token":      s.Token,
	}
}
func (s *SuccessResponse) SuccesResult() any {
	return map[string]any{
		"statuscode": s.StatusCode,
		"result":     s.Result,
		"token":      s.Token,
	}
}

// ErrorResponse adalah struktur respons kesalahan
type ErrorResponse struct {
	StatusCode int    `json:"statuscode"`
	Error      string `json:"error"`
	Detail     string `json:"detail"`
}

func (e *ErrorResponse) ErrorResult() any {
	return map[string]any{
		"statuscode": e.StatusCode,
		"err":        e.Error,
		"detail":     e.Detail,
	}
}
func (e *ErrorResponse) ErrorResultDetail() any {
	return map[string]any{
		"statuscode": e.StatusCode,
		"err":        e.Error,
		"detail":     e.Detail,
	}
}
