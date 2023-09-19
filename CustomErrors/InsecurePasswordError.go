package CustomErrors

type InsecurePassword struct {
	message string
}

func (taken InsecurePassword) Error() string {
	return "password needs to be at least 6 characters long, contain 1 upper case letter 1 lower case and a number"
}
