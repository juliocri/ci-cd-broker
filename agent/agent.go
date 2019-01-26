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
	Action string                 `json:"action"`
	Body   map[string]interface{} `json:"body"`
}

// Response defintion of response messages.
type Response struct {
	Status int                    `json:"status"`
	Body   map[string]interface{} `json:"body"`
}

// List of response-status codes.
const (
	statusOk       = 200
	statusError    = 500
	statusUAuth    = 401
	statusNotFound = 404
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
	List() Response
}
