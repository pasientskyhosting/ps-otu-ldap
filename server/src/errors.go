package main

// ErrorAPI - creates
type ErrorAPI struct {
	Error      ErrorAPIContent `json:"error"`
	StatusCode int             `json:"status_code"`
}

// ErrorAPIContent - creates
type ErrorAPIContent struct {
	Msg []ErrorAPIMessage `json:"messages"`
}

// ErrorAPIMessage - creates
type ErrorAPIMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ErrorAssetAlreadyExists - desc
func ErrorAssetAlreadyExists() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 409
	return a
}

// ErrorForbidden - desc
func ErrorForbidden() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 403
	a.AddMessage("ldap-users", "Forbidden")
	return a
}

// ErrorLDAPConn - desc
func ErrorLDAPConn() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 500
	return a
}

// ErrorAssetNotFound - desc
func ErrorAssetNotFound() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 404
	return a
}

// ErrorValidation - desc
func ErrorValidation() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 400
	return a
}

// AddMessage - desc
func (a *ErrorAPI) AddMessage(key string, val string) *ErrorAPI {
	a.Error.Msg = append(a.Error.Msg, ErrorAPIMessage{Key: key, Value: val})
	return a
}

// GetMessages - desc
func (a *ErrorAPI) GetMessages() []ErrorAPIMessage {
	return a.Error.Msg
}
