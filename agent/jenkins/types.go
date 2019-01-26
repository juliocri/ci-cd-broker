package jenkins

// CreateRequest definition fields to create a project in jenkins.
type CreateRequest struct {
	Name        string
	Description string
}

// DeleteRequest definition fields to delete a project in jenkins.
type DeleteRequest struct {
	Name string
}
