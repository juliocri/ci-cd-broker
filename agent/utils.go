package agent

import (
	"errors"
	"github.com/google/uuid"
)

// Error messages.
const (
	errIDe = "ID must to specified and be a valid uuid"
	errIDi = "ID is not a valid uuid"
	actErr = "Action field can not be skipped"
	bodErr = "Body field must to be specified"
)

//IsValid returns an error if an agent.Request is not valid.
func (r *Request) IsValid() error {

	if r.ID == "" {
		return errors.New(errIDe)
	}

	_, err := uuid.Parse(r.ID)
	if err != nil {
		return errors.New(errIDi)
	}

	if r.Action == "" {
		return errors.New(actErr)
	}

	if r.Body == nil {
		return errors.New(bodErr)
	}
	return nil
}
