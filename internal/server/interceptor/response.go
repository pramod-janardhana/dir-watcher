package interceptor

const (
	DEFAULT_ERROR_MSG         = "Something went wrong"
	DEFAULT_HTTP_ERROR_CODE   = 500
	DEFAULT_HTTP_SUCCESS_CODE = 200
)

// response imprements Response and defines the structure of the response to send.
type response struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	ErrMessage string      `json:"errMessage"`
}

// CreateResponse creates a new istence of response using the input params and returns.
// If errMsg is not provided then "default error message" will be used.
func NewResponse(success bool, data interface{}, errMsg string) *response {
	if !success && errMsg == "" {
		errMsg = DEFAULT_ERROR_MSG
	}
	return &response{
		Success:    success,
		Data:       data,
		ErrMessage: errMsg,
	}
}
