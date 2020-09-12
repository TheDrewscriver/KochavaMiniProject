package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)


type JSON struct {
	Endpoint struct {
		Method string `json:"method"`
		URL    string `json:"url"`
	} `json:"endpoint"`
	Data map[string]string `json:"data"`
}

func doReplace( URL string, DataMap map[string]string) string {
	//For each key value in the map
	for key, value:= range DataMap {
		//Surround the key with parantheses
		var StringKey string="{"+key+"}"
		//Replace the key with the value in the URL
		URL=strings.ReplaceAll(URL,StringKey,value)
	}
	//Replace any unchanged values with blanks
	regex := regexp.MustCompile(`={.*}`)
	URL= string(regex.ReplaceAll([]byte(URL),[]byte("=")))
	//Return the modified URL
	return URL
}

func main() {
	//Setup out
	log.SetOutput(os.Stdout)

	//Set up the consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		//Terminate the process if we cant create a consumer
		log.Print("ERROR: COULD NOT CREATE CONSUMER:"+err.Error())
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
			var messageReceived []byte =msg.Value
			log.Print("INFO: Message Received="+string(messageReceived))
			//Parse the incoming JSON into a an object
			incomingData := &JSON{}
			err2 := json.Unmarshal(messageReceived, &incomingData)
			if err2 != nil {
				//Print any errors
				log.Print("ERROR:Unable to parse JSON"+err2.Error())
			} else{
					//Replace all parameter placeholders with the actual value
					var URL string=doReplace(incomingData.Endpoint.URL,incomingData.Data)
					log.Print("INFO: Url being used:"+URL);
					log.Print("INFO: Method being used:"+incomingData.Endpoint.Method)
					//If Method is get/post
					if strings.EqualFold(incomingData.Endpoint.Method,"get"){
						resp, httpError :=http.Get(URL)
						if httpError != nil {
							log.Print("ERROR:Unable to send GET message due to error: "+ httpError.Error())
						}else if resp.StatusCode>400 {
							log.Print("ERROR:Unable to send message. Status code:"+string(resp.StatusCode))
						}
					}  else if strings.EqualFold(incomingData.Endpoint.Method,"post"){
					resp, httpError :=http.Post(URL,"*/*/",bytes.NewBuffer([]byte{}))
					if httpError != nil {
						log.Print("ERROR:Unable to send POST message due to error: "+ httpError.Error())
					}else if resp.StatusCode>400 {
						log.Print("ERROR:Unable to send message. Status code:"+string(resp.StatusCode))
					}
				}
			}
		} else {
			// The client will automatically try to recover from all errors.
			log.Print("ERROR: Consumer read error: "+ err.Error())
		}
	}

	c.Close()
}
