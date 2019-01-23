package agent

import (
	"fmt"
	"log"

	"github.com/bndr/gojenkins"
	"github.com/mitchellh/mapstructure"
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

// Create create something in jenkins.
func (j *Jenkins) Create(req Request) (Response, error) {
	// TODO implement the real statement.
	var resp Response
	var reqBody CreateBodyRequest

	resbody := map[string]interface{}{"msg": "Error"}
	err := mapstructure.Decode(req.Body, &reqBody)
	if err != nil {
		resp = Response{500, resbody}
		return resp, err
	}

	log.Printf("Trying to create a jenkins folder.")
	folder, err := j.Client.CreateFolder(reqBody.Name)
	if err != nil {
		resp = Response{500, resbody}
		return resp, err
	}
	log.Printf("Folder created successfuly: %v.", folder)

	resbody = map[string]interface{}{"msg": "Success"}
	resp = Response{200, resbody}
	return resp, err
}
