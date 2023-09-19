package CustomErrors

type UserTakenError struct {
	message string
}

func (taken UserTakenError) Error() string {
	return "User with given email is already registered"
}
