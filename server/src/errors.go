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

func ErrorAssetAlreadyExists() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 409
	return a
}

func ErrorForbidden() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 403
	a.AddMessage("ldap-users", "Forbidden")
	return a
}

func ErrorLDAPConn() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 500
	return a
}

func ErrorAssetNotFound() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 404
	return a
}

// ValidationError - creates
func ErrorValidation() ErrorAPI {
	a := ErrorAPI{}
	a.StatusCode = 400
	return a
}

func (a *ErrorAPI) AddMessage(key string, val string) *ErrorAPI {
	a.Error.Msg = append(a.Error.Msg, ErrorAPIMessage{Key: key, Value: val})
	return a
}

func (a *ErrorAPI) GetMessages() []ErrorAPIMessage {
	return a.Error.Msg
}
