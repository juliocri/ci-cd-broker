package jenkins

import (
	"errors"
)

// Error messages.
const (
	errName = "Name field must to be specified"
	errDesc = "Description field must to be specified"
)

// IsValid returns an error if a jenkins.CreateRequest is not valid.
func (cr *CreateRequest) IsValid() error {
	if cr.Name == "" {
		return errors.New(errName)
	}
	if cr.Description == "" {
		return errors.New(errDesc)
	}
	return nil
}

// IsValid return an error if a jenkins.UpdateRequest is not valid.
func (ur *UpdateRequest) IsValid() error {
	if ur.Description == "" {
		return errors.New(errDesc)
	}
	return nil
}

// IsValid return an error if a jenkins.DeleteRequest is not valid.
func (dr *DeleteRequest) IsValid() error {
	if dr.Name == "" {
		return errors.New(errName)
	}
	return nil
}
