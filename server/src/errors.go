package main

// APIError - creates
type APIError struct {
	Error      APIErrorContent `json:"error"`
	StatusCode int             `json:"status_code"`
}

// APIErrorContent - creates
type APIErrorContent struct {
	Msg []APIErrorMessage `json:"messages"`
}

// APIErrorMessage - creates
type APIErrorMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewAssetAlreadyExistsError - create
func NewAssetAlreadyExistsError() APIError {
	a := APIError{}
	a.StatusCode = 409
	return a
}

// NewForbiddenError - creates new forbidden error
func NewForbiddenError() APIError {
	a := APIError{}
	a.StatusCode = 403
	a.AddMessage("ldap-users", "Forbidden")
	return a
}

// LDAPConnError - creates
func LDAPConnError() APIError {
	a := APIError{}
	a.StatusCode = 500
	return a
}

// NewAssetNotFoundError - creates
func NewAssetNotFoundError() APIError {
	a := APIError{}
	a.StatusCode = 404
	return a
}

// NewValidationError - creates
func NewValidationError() APIError {
	a := APIError{}
	a.StatusCode = 400
	return a
}

// AddMessage - creates
func (a *APIError) AddMessage(key string, val string) *APIError {
	a.Error.Msg = append(a.Error.Msg, APIErrorMessage{Key: key, Value: val})
	return a
}

// GetMessages - creates
func (a *APIError) GetMessages() []APIErrorMessage {
	return a.Error.Msg
}
