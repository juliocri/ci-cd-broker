package agent

// Config CI/CD tool config
type Config struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Request defintion of request messages.
type Request struct {
	ID     string                 `json:"id"` // A uuid to identify the message.
	Action string                 `json:"action"`
	Body   map[string]interface{} `json:"body"`
}

// Response defintion of response messages.
type Response struct {
	ID     string                 `json:"id"` // A uuid received from a request.
	Status int                    `json:"status"`
	Body   map[string]interface{} `json:"body"`
}

// List of response-status codes.
const (
	StatusOk       = 200
	StatusError    = 500
	StatusUAuth    = 401
	StatusNotFound = 404
)

// Agent defines all methods to implement in Agent instances.
type Agent interface {
	// Get an agent connected to the proper vendor.
	Connect(Config) error
	// Create creates a  project in the vendor.
	Create(Request) (Response, error)
	// Delete deletes a project in the vendor.
	Delete(Request) (Response, error)
	// List returns a list with all projects in the vendor.
	List(Request) Response
	// Update modifies a project in the vendor.
	Update(Request) (Response, error)
}
