package users

type getUserError struct {
	desc string
}

func (e *getUserError) Error() string {
	return e.desc
}

type createUserError struct {
	desc string
}

func (e *createUserError) Error() string {
	return e.desc
}

type updateUserError struct {
	desc string
}

func (e *updateUserError) Error() string {
	return e.desc
}
