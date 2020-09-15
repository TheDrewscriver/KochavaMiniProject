package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type JSON struct {
	Endpoint struct {
		Method string `json:"method"`
		URL    string `json:"url"`
	} `json:"endpoint"`
	Data map[string]string `json:"data"`
}

func doReplace(URL string, DataMap map[string]string) string {
	//For each key value in the map
	for key, value := range DataMap {
		//Surround the key with parantheses
		var StringKey string = "{" + key + "}"
		//Replace the key with the value in the URL
		URL = strings.ReplaceAll(URL, StringKey, value)
	}
	//Replace any unchanged values with blanks
	regex := regexp.MustCompile(`={.*}`)
	URL = string(regex.ReplaceAll([]byte(URL), []byte("=")))
	//Return the modified URL
	return URL
}

func doSend(URL string, Method string) {
	var CurrentTime=time.Now()
	if strings.EqualFold(Method,"get") {
		resp, httpError := http.Get(URL)
		if httpError != nil {
			log.Error("Unable to send GET message due to error: " + httpError.Error())
		} else if resp.StatusCode >= http.StatusBadRequest {
			log.Error("Unable to send message. Status code:" + string(resp.StatusCode))
		}else {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			log.Info("Http Response body: "+bodyString)
		}
	}else if strings.EqualFold(Method, "post") {
		resp, httpError := http.Post(URL, "*/*/", bytes.NewBuffer([]byte{}))
		if httpError != nil {
			log.Error("Unable to send POST message due to error: " + httpError.Error())
		} else if  resp.StatusCode >= http.StatusBadRequest {
			log.Error("Unable to send message. Status code:" + string(resp.StatusCode))
		}else {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			log.Info("Http Response body: "+bodyString)
		}

	}
	var Diff=time.Since(CurrentTime)
	log.Info("Time for reponse: "+ strconv.FormatInt(Diff.Milliseconds(), 10)+" ms")
}



func main() {
	//Setup logger
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	//Set up the consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		//Terminate the process if we cant create a consumer
		log.Error("COULD NOT CREATE CONSUMER:" + err.Error())
		panic(err)
	}

	//Subscribe to the topic
	c.SubscribeTopics([]string{"kochavaPostback", "^aRegex.*[Tt]opic"}, nil)
	//While messages are still coming in
	for {
		//Read the message
		msg, readError := c.ReadMessage(-1)
		//If we don't have an error while reading
		if readError == nil {
			//Get the incoming message
			var messageReceived []byte = msg.Value
			log.Info("Message Received=" + string(messageReceived))
			//Parse the incoming JSON into a an object
			incomingData := &JSON{}
			err2 := json.Unmarshal(messageReceived, &incomingData)
			if err2 != nil {
				//Print any errors
				log.Error("Unable to parse JSON" + err2.Error())
			} else {
				//Replace all parameter placeholders with the actual value
				var URL string = doReplace(incomingData.Endpoint.URL, incomingData.Data)
				log.Info("Url being used:" + URL)
				log.Info("Method being used:" + incomingData.Endpoint.Method)
				doSend(URL,incomingData.Endpoint.Method)
			}
		} else {
			// The client will automatically try to recover from all errors.
			log.Error("Consumer read error: " + err.Error())
		}
	}

	c.Close()
}
