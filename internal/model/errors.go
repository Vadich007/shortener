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

type LinkDeletedError struct {
	ShortedLink string
}

func NewLinkDeletedError(shortedLink string) error {
	return &LinkDeletedError{ShortedLink: shortedLink}
}

func (err *LinkDeletedError) Error() string {
	return fmt.Sprintf("link is deleted (shorted_link = %s)", err.ShortedLink)
}
