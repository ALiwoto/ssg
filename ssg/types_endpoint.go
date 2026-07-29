package ssg

// EndpointResponse is the generalized form of a response from a HTTP API.
//
//	T field is already a pointer in this struct, please avoid passing a pointer
//
// type, to prevent from `Result` field being a pointer to a pointer.
type EndpointResponse[T any] struct {
	Success bool           `json:"success"`
	Result  *T             `json:"result"`
	Error   *EndpointError `json:"error"`
}

// EndpointError is the generalized form of an error from a HTTP API.
// this type should be used as a pointer in EndpointResponse.
type EndpointError struct {
	ErrorCode int    `json:"code"`
	ErrorType int    `json:"error_type"`
	Message   string `json:"message"`
	Origin    string `json:"origin"`
	Date      string `json:"date"`
}
