package agent

// Config CI/CD tool config
type Config struct {
	Name     string
	Host     string
	Port     int
	Username string
	Password string
}

// Agent defines all methods to implement in Agent instances.
type Agent interface {
	// Get an agent connected to the proper vendor.
	Connect(Config) error
}
