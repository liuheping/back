package context

import "errors"

const (
	PostMethodSupported = "only post method is allowed"
	CredentialsError    = "credentials error"
	TokenError          = "token error"
	UnauthorizedAccess  = "unauthorized access,user name or password error"
	NoRecord            = "no record"
	InvalidParam        = "invalid param"
)

var ErrNoRecord = errors.New(NoRecord)
var ErrUnAC = errors.New(UnauthorizedAccess)
var ErrInvalidParam = errors.New(InvalidParam)
