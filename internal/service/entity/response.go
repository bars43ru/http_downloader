package entity

import (
	"fmt"
)

type Response struct {
	Url        string
	StatusCode int
	Status     string
	Err        error
}

func (r Response) String() string  {
	if r.Err != nil {
		return fmt.Sprintf("url: %s, err: %v ", r.Url, r.Err)
	}
	return fmt.Sprintf("url: %s, status: %d, text: %s", r.Url, r.StatusCode, r.Status)
}
