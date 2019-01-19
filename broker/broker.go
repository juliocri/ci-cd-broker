package broker

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/optiopay/kafka"
	//"github.com/optiopay/kafka/proto"
)

// List of CI/CD supported.
var vendors = []string{"jenkins"}

// Default configuration from config file.
var config *Configuration

// Broker This struct stores the client reference.
type Broker struct {
	Client *kafka.Broker
}

// Configuration have all the requiered params to run the broker.
type Configuration struct {
	Host    string   `yaml:"host"`
	Port    int      `yaml:"port"`
	Vendors []Vendor `yaml:"vendors"`
}

// Vendor CI/CD tool config
type Vendor struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// setBrokerConfigFromFile reads a config file and parse the values into
// the structs.
func setBrokerConfigFromFile() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("%v ", err)
	}
	log.Printf("YAML config file was read.")

	var conf Configuration
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("%v ", err)
	}
	log.Printf("Config was parsed with no errors.")

	config = &conf
}

// getBrokerAdress get host and port from all the sources allowed.
func getBrokerAdress() []string {
	// Fisrt attemp go get adress values from the config.yaml file.
	setBrokerConfigFromFile()

	// If host and port are not set in the config file then defaulting
	if config.Host == "" {
		config.Host = "localhost"
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
	log.Printf("Broker address set successfuly at %v", address)

	return []string{address}
}

// GetBroker dial to kafka server and get the client to use.
func GetBroker() *Broker {
	// Set config to connect into kafka server.
	brokerAddress := getBrokerAdress()
	conf := kafka.NewBrokerConf("CI/CD Broker")
	conf.AllowTopicCreation = false
	// Connect to kafka cluster.
	broker, err := kafka.Dial(brokerAddress, conf)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer broker.Close()
	log.Printf("CI/CD Broker connection stablished at %v", brokerAddress)

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

// RunConsumers read messages from kafka in all ci/cd vendors topics and process
// them.
func (broker *Broker) RunConsumers() {
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
		go func(client kafka.Client, vendor Vendor) {
			// Consumers are intended to listen forever, but if for some reason or,
			// error the go rutine ends, we notify the sync group the routine is done.
			defer wg.Done()
			// Start creatin the consumer in the proper channel.
			inTopic := fmt.Sprintf("%v.requests", vendor.Name)
			outTopic := fmt.Sprintf("%v.responses", vendor.Name)
			partition := int32(0)
			conf := kafka.NewConsumerConf(inTopic, partition)
			conf.StartOffset = kafka.StartOffsetNewest

			log.Printf("Trying to create a consumer for topic `%v`", inTopic)
			consumer, err := client.Consumer(conf)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			log.Printf("Consumer created successfuly for topic `%v`", inTopic)
			// Infinite loop to listen the message queue in the proper topic.
			for {
				msg, err := consumer.Consume()
				if err != nil {
					if err != kafka.ErrNoData {
						log.Printf("Error: %v", err)
					}
					break
				}
				log.Printf("Message in topic `%v` was read.", inTopic)

				// TODO : process message and redirect to the proper ci/cd api.
				// Creating a producer to write in the vendor topic a response.
				log.Printf("Trying to push a message in topic `%v`", outTopic)
				producer := client.Producer(kafka.NewProducerConf())
				msg.Value = append(msg.Value, "bblb"...)
				if _, err := producer.Produce(outTopic, partition, msg); err != nil {
					log.Fatalf("%v", err)
				}

				log.Printf("Message in topic `%v` was written.", outTopic)
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
