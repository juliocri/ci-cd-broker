package agent

// MockAgent implemenation of agent for unit testing.
type MockAgent struct {
	Client string
}

//Connect -
func (m *MockAgent) Connect(conf Config) error {
	m.Client = conf.Name
	return nil
}

// Create -
func (m *MockAgent) Create(req Request) (Response, error) {
	return Response{req.ID, StatusOk, map[string]interface{}{
		"msg":     req.Action,
		"details": "ok",
	}}, nil
}

// Delete -
func (m *MockAgent) Delete(req Request) (Response, error) {
	return Response{req.ID, StatusOk, map[string]interface{}{
		"msg":     req.Action,
		"details": "ok",
	}}, nil
}

// List -
func (m *MockAgent) List(req Request) Response {
	var list []map[string]interface{}
	return Response{req.ID, StatusOk, map[string]interface{}{
		"msg":      req.Action,
		"details":  "ok",
		"projects": list,
	}}
}

// Update -
func (m *MockAgent) Update(req Request) (Response, error) {
	return Response{req.ID, StatusOk, map[string]interface{}{
		"msg":     req.Action,
		"details": "ok",
	}}, nil
}
