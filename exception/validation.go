package exception

type (
	WrongCredentialError struct {
		Message string
	}
)

func (e WrongCredentialError) Error() string {
	return e.Message
}
