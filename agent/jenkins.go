package agent

import (
	"fmt"
	"log"

	"github.com/bndr/gojenkins"
	"github.com/mitchellh/mapstructure"
	"github.com/juliocri/ci-cd-broker/agent/jenkins"
	"github.com/juliocri/ci-cd-broker/agent/jenkins/configs"
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
	// Decoding the request in an specific create request for jenkins.
	err = mapstructure.Decode(req.Body, &request)
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}
	// Verifying if all mandatory fields are set in the request.
	err = request.IsValid()
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}

	// TODO build the config string with all the parameters recieved by the GUI.
	log.Printf("Trying to create project `%v`", request.Name)
	configString := fmt.Sprintf(configs.XMLCreate, request.Description)
	job, err := j.Client.CreateJob(configString, request.Name)
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}
	d := fmt.Sprintf("Project `%v` has been created.", job.GetName())
	log.Printf(d)

	// Invoking job in order to scan repo for first time.
	log.Printf("Starting a quick self-scan in project `%v`.", job.GetName())
	var params map[string]string
	_, _ = job.InvokeSimple(params)
	log.Printf("Scan in project `%v` executed successfuly.", job.GetName())

	// Finally a successfull response.
	response = Response{
		req.ID,
		StatusOk,
		map[string]interface{}{"msg": "success", "details": d},
	}

	return response, err
}

// Delete deletes a job in jenkins.
func (j *Jenkins) Delete(req Request) (Response, error) {
	var response Response
	var request jenkins.DeleteRequest
	var err error
	// Decoding request into a delete request structure for jenkins.
	err = mapstructure.Decode(req.Body, &request)
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}
	// Verifying if all mandatory fields are set in the request.
	err = request.IsValid()
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}

	log.Printf("Trying to delete project `%v`.", request.Name)
	_, err = j.Client.DeleteJob(request.Name)
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err,
		}}
		return response, err
	}

	d := fmt.Sprintf("Project `%v` has been delete.", request.Name)
	log.Printf(d)
	response = Response{
		req.ID,
		StatusOk,
		map[string]interface{}{"msg": "success", "details": d},
	}

	return response, err
}

// List return a list with all pipelines in jenkins.
func (j *Jenkins) List(req Request) Response {
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
	d := "List of projects retrieved."
	log.Printf(d)
	response.ID = req.ID
	response.Status = StatusOk
	response.Body = map[string]interface{}{
		"msg":      "success",
		"details":  d,
		"projects": list,
	}

	return response
}

// Update updates a pipeline in jenkins.
func (j *Jenkins) Update(req Request) (Response, error) {
	var response Response
	var request jenkins.UpdateRequest
	var err error
	// Decoding the request's body into a update request for jenkins.
	err = mapstructure.Decode(req.Body, &request)
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}
	// Verifying if all mandatory fields are set in the update request.
	err = request.IsValid()
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}
	// First step to update is get the job information.
	job, err := j.Client.GetJob(request.Name)
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}
	// TODO update or add a new var for the config string to update, instead of
	// using the create config string.
	log.Printf("Trying to update project `%v`.", job.GetName())
	configString := fmt.Sprintf(configs.XMLCreate, request.Description)
	err = job.UpdateConfig(configString)
	if err != nil {
		response = Response{req.ID, StatusError, map[string]interface{}{
			"msg":     "error",
			"details": err.Error(),
		}}
		return response, err
	}
	// Update complete.
	d := fmt.Sprintf("Project `%v` has been updated.", job.GetName())
	log.Printf(d)
	response.ID = req.ID
	response.Status = StatusOk
	response.Body = map[string]interface{}{"msg": "success", "details": d}

	return response, err
}
