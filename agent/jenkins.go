package agent

import (
	"fmt"
	"log"

	"github.com/bndr/gojenkins"
	"github.com/mitchellh/mapstructure"
	"github.intel.com/kubernetes/ci-cd-broker/agent/jenkins"
)

// Jenkins implemenation of agent.
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

// Create creates a job in jenkins.
func (j *Jenkins) Create(req Request) (Response, error) {
	var response Response
	var request jenkins.CreateRequest
	var err error

	err = mapstructure.Decode(req.Body, &request)
	if err != nil {
		response = Response{statusError, map[string]interface{}{"msg": err}}
		return response, err
	}

	// TODO build the config string with all the parameters recieved by the GUI.
	log.Printf("Trying to create project `%v`", request.Name)
	configString := fmt.Sprintf(jenkins.XMLCreate, request.Description)
	job, err := j.Client.CreateJob(configString, request.Name)
	if err != nil {
		response = Response{statusError, map[string]interface{}{"msg": err}}
		return response, err
	}
	log.Printf("Project `%v` has been created.", job.GetName())

	// Invoking job in order to scan repo for first time.
	log.Printf("Starting a quick self-scan in project `%v`.", job.GetName())
	var params map[string]string
	_, _ = job.InvokeSimple(params)
	log.Printf("Scan in project `%v` executed successfuly.", job.GetName())

	// Finally a successfull response.
	response = Response{
		statusOk,
		map[string]interface{}{"msg": "Success!"},
	}

	return response, err
}

// Delete deletes a job in jenkins.
func (j *Jenkins) Delete(req Request) (Response, error) {
	var response Response
	var request jenkins.CreateRequest
	var err error

	err = mapstructure.Decode(req.Body, &request)
	if err != nil {
		response = Response{statusError, map[string]interface{}{"msg": err}}
		return response, err
	}

	log.Printf("Trying to delete project `%v`", request.Name)
	_, err = j.Client.DeleteJob(request.Name)
	if err != nil {
		response = Response{statusError, map[string]interface{}{"msg": err}}
		return response, err
	}
	log.Printf("Project `%v` has been delete.", request.Name)

	return response, err
}

// List return a list with all pipelines in jenkins.
func (j *Jenkins) List() Response {
	var response Response
	var list []map[string]interface{}

	log.Printf("Trying to get the list of projects in jenkins.")
	jobs, _ := j.Client.GetAllJobs()
	for _, job := range jobs {
		item := map[string]interface{}{
			"name":        job.GetName(),
			"description": job.GetDescription(),
		}
		list = append(list, item)
	}
	log.Printf("List of projects retrieved.")

	response.Status = statusOk
	response.Body = map[string]interface{}{
		"msg":  "Success!",
		"jobs": list,
	}

	return response
}
