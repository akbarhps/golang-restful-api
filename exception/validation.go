package exception

type (
	InvalidCredentialError struct {
		Message string
	}
	InvalidSignatureError struct {
		Message string
	}
)

func (e InvalidCredentialError) Error() string {
	return e.Message
}

func (e InvalidSignatureError) Error() string {
	return e.Message
}
