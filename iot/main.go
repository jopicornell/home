package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yryz/ds18b20"
	"log"
	"net/http"
)

var sensorMap = map[string]string {
	"28-012032eb3734": "nick-busk",
	"28-012032cd28ce": "office",
}

func main() {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	fmt.Printf("sensor IDs: %v\n", sensors)
	var request []map[string]interface{}
	for _, sensor := range sensors {
		t, err := ds18b20.Temperature(sensor)
		if err != nil {
			log.Printf("Error getting temps: %s", err.Error())
			continue
		}
		fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor, t)
		request = append(request, map[string]interface{}{
			"temperature": t,
			"type": sensorMap[sensor],
		})
	}
	body, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	response, err := http.Post("https://umczz0pvpc.execute-api.eu-central-1.amazonaws.com/prod/temperatures", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created temperatures. Response: %+v", response)
}