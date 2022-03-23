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

type Status struct {
	HeaterOn    bool
	LightOn     bool
	BaskingTemp float64
	ColdTemp    float64
}

var currentStatus Status = Status{
	HeaterOn:    false,
	LightOn:     false,
	BaskingTemp: 0,
	ColdTemp:    0,
}

func main() {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}
	fmt.Printf("sensor IDs: %v\n", sensors)
	var request []map[string]interface{}
	currentBaskingTemp := getBaskingTemp()
	condition := getCurrentCondition()
	defaultGpioValue := 0
	if condition.IsHeaterOn(time.Now(), float32(currentBaskingTemp)) {
		defaultGpioValue = 1
	}
	gpioHeater, err := gpiod.RequestLine("gpiochip0", 14, gpiod.AsOutput(defaultGpioValue))
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
			request = append(request, map[string]interface{}{
				"temperature": temperature,
				"type":        sensorMap[sensor],
			})
			if sensorMap[sensor] == "nick-basking" {
				currentStatus.BaskingTemp = temperature
				condition := getCurrentCondition()
				if !currentStatus.HeaterOn && condition.IsHeaterOn(time.Now(), float32(temperature)) {
					fmt.Printf("Activating %s temperature reached %.2f less than %.2f", sensorMap[sensor], temperature, condition.IdealTemperature(time.Now()))
					gpioHeater.SetValue(1)
				} else if currentStatus.HeaterOn && !condition.IsHeaterOn(time.Now(), float32(temperature)) {
					fmt.Printf("Deactivating %s temperature reached %.2f less than %.2f", sensorMap[sensor], temperature, condition.IdealTemperature(time.Now()))
					gpioHeater.SetValue(0)
				}
				currentStatus.HeaterOn = condition.IsHeaterOn(time.Now(), float32(currentStatus.BaskingTemp))
				currentStatus.LightOn = condition.IsLightOn(time.Now())
			}
			if sensorMap[sensor] == "nick-cold" {
				currentStatus.ColdTemp = temperature
			}
		}
		time.Sleep(5 * time.Second)
	}
	//	response, err := http.Post("https://umczz0pvpc.execute-api.eu-central-1.amazonaws.com/prod/temperatures", "application/json", bytes.NewBuffer(body))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("Created temperatures. Response: %+v", response)
}

func getBaskingTemp() float64 {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}
	for _, sensor := range sensors {
		temperature, err := ds18b20.Temperature(sensor)
		if err != nil {
			log.Printf("Error getting temps: %s", err.Error())
			return 21
		}
		if sensorMap[sensor] == "nick-basking" {
			return temperature
		}
	}
	return 21
}

func getCurrentCondition() models.Condition {
	return &models.SunsetSunriseCondition{
		MinimumHours:     9,
		MaximumHours:     12,
		TemperatureNight: 21,
		TemperatureDay:   31,
	}
}

func getHeaterLine() (*gpiod.Line, error) {
	return gpiod.RequestLine("gpiochip0", 14, gpiod.AsOutput(1))
}
