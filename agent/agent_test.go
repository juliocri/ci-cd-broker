package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mock MockAgent

func TestConnect(t *testing.T) {
	conf := Config{"MockVendor", "mock-agent.com", 8080, "mockuser", "mockpasswd"}
	res := mock.Connect(conf)
	assert := assert.New(t)
	assert.Equal(conf.Name, mock.Client, "Mock client not nil.")
	assert.Equal(res, nil, "Error is not expected.")
}

func TestCreate(t *testing.T) {
	req := Request{"id-123", "create", nil}
	res, err := mock.Create(req)
	assert := assert.New(t)
	assert.Equal(StatusOk, res.Status, "Response status 200 expected.")
	assert.Equal(err, nil, "Error is not expected.")
}

func TestDelete(t *testing.T) {
	req := Request{"id-123", "delete", nil}
	res, err := mock.Delete(req)
	assert := assert.New(t)
	assert.Equal(StatusOk, res.Status, "Response status 200 expected.")
	assert.Equal(err, nil, "Error is not expected.")
}

func TestUpdate(t *testing.T) {
	req := Request{"id-123", "update", nil}
	res, err := mock.Update(req)
	assert := assert.New(t)
	assert.Equal(StatusOk, res.Status, "Response status 200 expected.")
	assert.Equal(err, nil, "Error is not expected.")
}

func TestList(t *testing.T) {
	req := Request{"id-123", "list", nil}
	res := mock.List(req)
	var exp []map[string]interface{}
	assert := assert.New(t)
	assert.Equal(StatusOk, res.Status, "Response status 200 expected.")
	assert.Equal(res.Body["projects"], exp, "An empty list of projects is expected.")
}
