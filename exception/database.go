package exception

type (
	RecordDuplicateError struct {
		Message string
	}

	RecordNotFoundError struct {
		Message string
	}
)

func (e RecordDuplicateError) Error() string {
	return e.Message
}

func (e RecordNotFoundError) Error() string {
	return e.Message
}
