package model

import "fmt"

type LinkAlreadyExistError struct {
	Err         error
	ShortedLink string
}

func NewLinkAlreadyExistError(shortedLink string) error {
	return &LinkAlreadyExistError{
		ShortedLink: shortedLink,
		Err:         fmt.Errorf("record already exist (shorted_link = %s)", shortedLink),
	}
}

func (err *LinkAlreadyExistError) Error() string {
	return err.Err.Error()
}
