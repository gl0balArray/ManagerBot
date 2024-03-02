package error

type DataNotFoundError struct{
	What string
}

func (e DataNotFoundError) Error() string{
	return e.What
}