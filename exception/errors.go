package exception

type (
	Errors struct {
		Errors []error `json:"errors"`
	}

	TokenError struct {
		Message string `json:"message"`
	}

	NotFoundError struct {
		Message string `json:"message"`
	}

	DatabaseError struct {
		Message string `json:"message"`
	}

	WrongPasswordError struct {
        Message string `json:"message"`
    }

	FieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

func (e TokenError) Error() string {
	return e.Message
}

func (e NotFoundError) Error() string {
	return e.Message
}

func (e DatabaseError) Error() string {
	return e.Message
}

func (e WrongPasswordError) Error() string {
    return e.Message
}

func (e FieldError) Error() string {
	return e.Message
}
