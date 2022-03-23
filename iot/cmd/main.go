package main

import (
	"fmt"
	"github.com/jopicornell/thermonick/models"
	"github.com/warthog618/gpiod"
	"github.com/yryz/ds18b20"
	"log"
	"time"
	//	"net/http"
)

var sensorMap = map[string]string{
	"28-012032eb3734": "nick-basking",
	"28-012032cd28ce": "nick-cold",
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
			temperature, err := ds18b20.Temperature(sensor)
			if err != nil {
				log.Printf("Error getting temps: %s", err.Error())
				continue
			}
			fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor, temperature)
			request = append(request, map[string]interface{}{
				"temperature": temperature,
				"type":        sensorMap[sensor],
			})
			if sensorMap[sensor] == "nick-busk" {
				condition := getCurrentCondition()
				if condition.IsHeaterOn(time.Now(), float32(temperature)) {
					fmt.Printf("Activating %s", sensorMap[sensor])
					gpioHeater.SetValue(1)
				} else {
					fmt.Printf("Deactivating %s", sensorMap[sensor])
					gpioHeater.SetValue(0)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
	//	response, err := http.Post("https://umczz0pvpc.execute-api.eu-central-1.amazonaws.com/prod/temperatures", "application/json", bytes.NewBuffer(body))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("Created temperatures. Response: %+v", response)
}

func getCurrentCondition() models.Condition {
	return &models.SunsetSunriseCondition{
		MinimumHours:     9,
		MaximumHours:     12,
		TemperatureNight: 21,
		TemperatureDay:   31,
	}
}

func getSensorLine(sensor string) (*gpiod.Line, error) {
	return gpiod.RequestLine("gpiochip0", 14, gpiod.AsOutput(1))
}
