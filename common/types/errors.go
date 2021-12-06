package types

// APIError is one of the following
// 400 Bad Request: The data of your request is invalid or incomplete.
// 401 Unauthorized: You are not correctly authenticated.
// 403 Forbidden: You are authenticated, but you are not authorized.
// 404 Not Found: The resource you try to access does not exist.
// 415 Unsupported Media Type: Invalid Content-Type or Accept HTTP header.
// 417 Expectation Failed: Your token quota has been exhausted.
// 422 Unprocessable Entity: Validation error (missing required attributes, number not in range etc.)
// 423 Locked: The API is locked for write operations during maintenance.
// 429 Too Many Requests: You have been rate limited.
// 500 Internal Server Error: Something wrong happened in the server
// 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout: Temporary communication failure.
// There may be additional error codes in certain circumstances. Please refer to the HTTP error code documentation for more information.

// APIError API Error
type APIError struct {
	StatusCode int
	error
}

// NewAPIError returns a new wrapped APIError
func NewAPIError(err error) *APIError {
	return &APIError{
		StatusCode: 500,
		error:      err,
	}
}
