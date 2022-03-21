package main

import (
//	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yryz/ds18b20"
	"log"
	"github.com/warthog618/gpiod"
	"time"
//	"net/http"
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

	gpioHeater, err := gpiod.RequestLine("gpiochip0", 14, gpiod.AsOutput(1))
	if err != nil {
		log.Fatal(err)
	}
	for {
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
			if sensorMap[sensor] == "nick-busk" {
				if (t < 21) {
					fmt.Printf("Activating %s", sensorMap[sensor])
					gpioHeater.SetValue(1)
				} else {
					fmt.Printf("Deactivating %s", sensorMap[sensor])
					gpioHeater.SetValue(0)
				}
			}
					}
		time.Sleep(5 * time.Second)
	}
	body, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("tempObj: %+v", body)
//	response, err := http.Post("https://umczz0pvpc.execute-api.eu-central-1.amazonaws.com/prod/temperatures", "application/json", bytes.NewBuffer(body))
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("Created temperatures. Response: %+v", response)
}
