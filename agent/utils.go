package agent

import (
	"errors"
)

// Error messages.
const (
	actErr = "Action field can not be skipped"
	bodErr = "Body field must to be specified"
)

//IsValid returns an error if an agent.Request is not valid.
func (r *Request) IsValid() error {
	if r.Action == "" {
		return errors.New(actErr)
	}
	if r.Body == nil {
		return errors.New(bodErr)
	}
	return nil
}
