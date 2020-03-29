package er

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (err StatusError) Status() int {
	return err.Code
}

func (err StatusError) Error() string {
	return err.Err.Error()
}
