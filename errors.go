package comusic

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrUnauthenticated   = Error("unauthenticated")
	ErrAuthProcessFailed = Error("authentication failed")
	ErrResourceNotFound  = Error("resource not found")
)
