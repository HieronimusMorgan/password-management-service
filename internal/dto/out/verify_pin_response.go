package out

type VerifyPinCodeResponse struct {
	ClientID  string `json:"client_id"`
	RequestID string `json:"request_id"`
	Valid     bool   `json:"valid"`
}
