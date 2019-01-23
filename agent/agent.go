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

// CreateBodyRequest definition of body fields for create requests.
type CreateBodyRequest struct {
	Name string
}

// Response defintion of response messages.
type Response struct {
	Status int                    `json:"status"`
	Body   map[string]interface{} `json:"body"`
}

// Agent defines all methods to implement in Agent instances.
type Agent interface {
	// Get an agent connected to the proper vendor.
	Connect(Config) error
	// TODO Create something in jenkins
	Create(Request) (Response, error)
}
