package agent

import (
	"fmt"
	"log"

	"github.com/bndr/gojenkins"
)

// Jenkins alv
type Jenkins struct {
	Client *gojenkins.Jenkins
}

// Connect returns a jenkins agent client.
func (j *Jenkins) Connect(conf Config) error {
	address := fmt.Sprintf("%v:%v", conf.Host, conf.Port)
	jenkins := gojenkins.CreateJenkins(
		nil,
		address,
		conf.Username,
		conf.Password,
	)
	// Provide CA certificate if server is using self-signed certificate
	// caCert, _ := ioutil.ReadFile("/tmp/ca.crt")
	// jenkins.Requester.CACert = caCert
	log.Printf("Agent jenkins is trying to connect to client.")
	_, err := jenkins.Init()
	if err != nil {
		return err
	}
	log.Printf("Agent jenkins is connected to client.")

	j.Client = jenkins
	return nil
}
