package broker

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/juliocri/ci-cd-broker/agent"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
)

// List of CI/CD supported.
var vendors = []string{"jenkins"}

// Default configuration from config file.
var config *Config

// Broker This struct stores the client reference.
type Broker struct {
	Client *kafka.Broker
}

// Config have all the requiered params to run the broker.
type Config struct {
	Host    string         `yaml:"host"`
	Port    int            `yaml:"port"`
	Vendors []agent.Config `yaml:"vendors"`
}

// setConfigFromFile reads a config file and parse the values into
// the structs.
func setConfigFromFile() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("%v ", err)
	}
	log.Printf("YAML config file was read.")

	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("%v ", err)
	}
	log.Printf("Config was parsed with no errors.")

	config = &conf
}

// getAdress get host and port from all the sources allowed.
func getAdress() []string {
	// Fisrt attemp go get adress values from the config.yaml file.
	setConfigFromFile()

	// If host and port are not set in the config file then defaulting
	if config.Host == "" {
		config.Host = "kafka"
	}
	if config.Port == 0 {
		config.Port = 9092
	}

	// Search if the values are present as command line args, if there are present
	// then values are overwritten otherwise yaml file or default values are set.
	var kafkaAddress = flag.String(
		"host",
		config.Host,
		"Kafka address/hostname string",
	)
	var kafkaPort = flag.Int(
		"port",
		config.Port,
		"Kafka port number",
	)
	flag.Parse()
	// Assign ultimate values to the broker info
	address := fmt.Sprintf("%v:%v", *kafkaAddress, *kafkaPort)
	log.Printf("Broker address set successfuly at `%v`.", address)

	return []string{address}
}

// Get dials to kafka server and get the client to use.
func Get() *Broker {
	// Set config to connect into kafka server.
	brokerAddress := getAdress()
	conf := kafka.NewBrokerConf("CI/CD Broker")
	conf.AllowTopicCreation = false
	// Connect to kafka cluster.
	broker, err := kafka.Dial(brokerAddress, conf)
	if err != nil {
		log.Fatalf("Error: %s.", err)
	}
	defer broker.Close()
	log.Printf("CI/CD Broker connection stablished at `%v`.", brokerAddress)

	return &Broker{broker}
}

// verifySupportForVendor returns true if a vendor name set in the config is,
// supported by the broker.
func verifySupportForVendor(name string) bool {
	var supported = false
	for _, vendor := range vendors {
		if vendor == name {
			supported = true
			break
		}
	}

	return supported
}

// getVendorAgent return an agent client accordingly to the vendor name.
func getVendorAgent(name string) (agent.Agent, error) {
	var agt agent.Agent
	var err error
	switch name {
	case "jenkins":
		agt = &agent.Jenkins{}

	default:
		msg := fmt.Sprintf("Vendor client `%v` is not supported.", name)
		err = errors.New(msg)
	}
	return agt, err
}

// fetch returns a json string as a response from a request.
func fetch(agt agent.Agent, msg []byte) ([]byte, error) {
	var req agent.Request
	var res agent.Response
	var err error
	var response []byte

	// template response for error output.
	tmp := `{"id":"%v","status":%v,"body":{"msg":"error","details":"%v"}}`
	// Parse json request to agent.Request type
	err = json.Unmarshal(msg, &req)
	if err != nil {
		r := fmt.Sprintf(tmp, "invalid-message", agent.StatusError, err)
		return []byte(r), err
	}

	err = req.IsValid()
	if err != nil {
		r := fmt.Sprintf(tmp, req.ID, agent.StatusError, err)
		return []byte(r), err
	}

	// Execute the action, and get an agent.Response
	switch req.Action {
	case "create":
		res, err = agt.Create(req)

	case "delete":
		res, err = agt.Delete(req)

	case "list":
		res, err = agt.List(req), nil

	case "update":
		res, err = agt.Update(req)

	default:
		msg := fmt.Sprintf("Action `%v` is not implemented in agent.", req.Action)
		err = errors.New(msg)
		r := fmt.Sprintf(tmp, req.ID, agent.StatusError, err)
		return []byte(r), err
	}

	// Try to Parse the agent.Response into a json stream.
	response, errd := json.Marshal(res)
	if errd != nil {
		r := fmt.Sprintf(tmp, req.ID, agent.StatusError, errd)
		return []byte(r), errd
	}

	// return the action response parsed to []byte.
	return response, err
}

// Run read messages from kafka in all ci/cd vendors topics and process
// them.
func (broker *Broker) Run() {
	var wg sync.WaitGroup
	// Create a consumer in each vendor supported.
	for _, vendor := range config.Vendors {
		// Verify if the vendor.Name exists in supported vendor list, otherwise the,
		// creation of consumer will be skipped.
		if support := verifySupportForVendor(vendor.Name); support != true {
			log.Printf(
				"Broker has not support for vendor `%v`, skipping...",
				vendor.Name,
			)
			continue
		}
		// We start to create the consumers in the current Topic, Creating a go
		// routine for each topic.
		wg.Add(1)
		go func(client kafka.Client, vendor agent.Config) {
			// Consumers are intended to listen forever, but if for some reason or,
			// error the go rutine ends, we notify the sync group the routine is done.
			defer wg.Done()
			// Start creatin the consumer in the proper channel.
			inTopic := fmt.Sprintf("%v-requests", vendor.Name)
			outTopic := fmt.Sprintf("%v-responses", vendor.Name)
			partition := int32(0)
			conf := kafka.NewConsumerConf(inTopic, partition)
			conf.StartOffset = kafka.StartOffsetNewest

			log.Printf("Trying to create a consumer for topic `%v`.", inTopic)
			consumer, err := client.Consumer(conf)
			if err != nil {
				log.Fatalf("Error: %v.", err)
			}
			log.Printf("Consumer created successfuly for topic `%v`.", inTopic)

			// Get the proper vendor client accordingly to the vendor name.
			agt, err := getVendorAgent(vendor.Name)
			if err != nil {
				log.Fatalf("%v", err)
			}

			// Make a connection with the vendor client.
			err = agt.Connect(vendor)
			if err != nil {
				log.Fatalf("%v", err)
			}

			log.Printf("Trying to create a producer for topic `%v`.", outTopic)
			producer := client.Producer(kafka.NewProducerConf())
			log.Printf("Producer created successfuly for topic `%v`.", outTopic)

			// Till this point we have the machinery required to read-process-write.
			// Infinite loop to listen the message queue in the proper topic.
			log.Printf("Consumer for topic `%v` is listening for messages.", inTopic)
			for {
				msg, err := consumer.Consume()
				if err != nil {
					if err != kafka.ErrNoData {
						log.Printf("Error: %v.", err)
					}
					break
				}
				log.Printf("A message in topic `%v` was read.", inTopic)

				// Async process of message.
				wg.Add(1)
				go func(a agent.Agent, m *proto.Message, p kafka.Producer) {
					defer wg.Done()
					// Build the msg and sending the action request to the proper client.
					resp, err := fetch(a, m.Value)
					m.Value = resp
					// We get an error response from an action message.
					if err != nil {
						log.Printf("Action response error: %v.", err)
						// If a error occurs we need to know it and write it in kafka.
						// Logging the error and write an error response in the topic.
						log.Printf("Trying to write an error message in topic `%v`.", outTopic)
						if _, err := p.Produce(outTopic, partition, m); err != nil {
							// Writting in kafka results in error.
							// Logging the error, but keep listening the topic.
							log.Printf("Warning: %v.", err)
							return
						}
						log.Printf("Error Message in topic `%v` was written.", outTopic)
						return
					}

					// We get a successfully response from an action message.
					log.Printf("Trying to write a success message in topic `%v`.", outTopic)
					if _, err := p.Produce(outTopic, partition, m); err != nil {
						// Writting in kafka results in error.
						// Logging the error, but keep listening the topic.
						log.Printf("Warning: %v.", err)
						return
					}
					log.Printf("Success Message in topic `%v` was written.", outTopic)
				}(agt, msg, producer)

			}
			// If getting this point then means the consumer ends for some reazon.
			log.Printf("Consumer for topic `%v` has ended.", inTopic)
			// End of gorutine.
		}(broker.Client, vendor)
	}

	log.Printf("CI/CD Broker starting operations...")
	// Listen all topics until all go rutines ends.
	wg.Wait()
	// If we recieve all gorutines done signal then end program.
	log.Printf("CI/CD Broker ending operations...")
}
