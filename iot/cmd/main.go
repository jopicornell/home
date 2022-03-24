package main

import (
	"bufio"
	"fmt"
	"github.com/jopicornell/thermonick/models"
	"github.com/warthog618/gpiod"
	"github.com/yryz/ds18b20"
	"log"
	"os"
	"time"
	//	"net/http"
)

type Status struct {
	HeaterOn    bool
	LightOn     bool
	BaskingTemp float64
	ColdTemp    float64
}

var sensorMap = map[string]string{
	"28-012032eb3734": "basking",
	"28-012032cd28ce": "cold",
}

var sensorNameMap = map[string]string{
	"basking": "28-012032eb3734",
	"cold":    "28-012032cd28ce",
}

var currentStatus Status = Status{
	HeaterOn:    false,
	LightOn:     false,
	BaskingTemp: 0,
	ColdTemp:    0,
}

var HeaterLine *gpiod.Line
var HeaterLineNumber int

func main() {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}
	fmt.Printf("sensor IDs: %v\n", sensors)
	currentBaskingTemp := getBaskingTemp()
	fmt.Printf("currentBaskingTemp: %f\n", currentBaskingTemp)
	ReadCommandsRoutine()
	for {
		baskingTemperature := getBaskingTemp()
		currentStatus.BaskingTemp = baskingTemperature
		if ShouldActivateHeater(baskingTemperature) {
			ActivateHeater(baskingTemperature)
		} else if ShouldDeactivateHeater(baskingTemperature) {
			DeactivateHeater(baskingTemperature)
		}
		if ShouldActivateLight() {
			currentStatus.LightOn = true
		} else if ShouldDeactivateLight() {
			currentStatus.LightOn = false
		}
		currentStatus.ColdTemp = getColdTemp()
		time.Sleep(5 * time.Second)
	}
	//	response, err := http.Post("https://umczz0pvpc.execute-api.eu-central-1.amazonaws.com/prod/temperatures", "application/json", bytes.NewBuffer(body))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("Created temperatures. Response: %+v", response)
}

func ReadCommandsRoutine() {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			cmd, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			if cmd == "status\n" {
				fmt.Printf("currentStatus heater %v %v %f %f \n", currentStatus.HeaterOn, currentStatus.LightOn, currentStatus.BaskingTemp, currentStatus.ColdTemp)
			}
		}

	}()
}

func DeactivateHeater(baskingTemperature float64) {
	fmt.Printf("Deactivating heater temperature reached %.2f more than %.2f\n", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getHeaterLine().SetValue(0)
	if err != nil {
		log.Printf("Error deactivating temps: %s", err.Error())
	}
	currentStatus.HeaterOn = false
}

func ActivateHeater(baskingTemperature float64) {
	fmt.Printf("Activating heater temperature reached %.2f less than %.2f\n", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getHeaterLine().SetValue(1)
	if err != nil {
		log.Printf("Error activating heater: %+v", err)
	}
	currentStatus.HeaterOn = true
}

func getBaskingTemp() float64 {
	return getTemperature(sensorNameMap["basking"])
}

func getColdTemp() float64 {
	return getTemperature(sensorNameMap["cold"])
}

func getTemperature(sensorName string) float64 {
	temperature, err := ds18b20.Temperature(sensorName)
	if err != nil {
		log.Printf("Error getting basking temperature: %s", err.Error())
		return 21
	}
	return temperature
}

func ShouldActivateHeater(temperature float64) bool {
	condition := getCurrentCondition()
	return !currentStatus.HeaterOn && condition.IsHeaterOn(time.Now(), temperature)
}

func ShouldDeactivateHeater(temperature float64) bool {
	condition := getCurrentCondition()
	return currentStatus.HeaterOn && !condition.IsHeaterOn(time.Now(), temperature)
}

func ShouldActivateLight() bool {
	condition := getCurrentCondition()
	return !currentStatus.HeaterOn && condition.IsLightOn(time.Now())
}

func ShouldDeactivateLight() bool {
	condition := getCurrentCondition()
	return currentStatus.HeaterOn && !condition.IsLightOn(time.Now())
}

func getCurrentCondition() models.Condition {
	return &models.SunsetSunriseCondition{
		MinimumHours:     9,
		MaximumHours:     12,
		TemperatureNight: 21,
		TemperatureDay:   31,
	}
}

func getDefaultGpioValue() int {
	condition := getCurrentCondition()
	if condition.IsHeaterOn(time.Now(), getBaskingTemp()) {
		return 1
	}
	return 0
}

func getHeaterLine() *gpiod.Line {
	if HeaterLine == nil {
		line, err := gpiod.RequestLine("gpiochip0", HeaterLineNumber, gpiod.AsOutput(getDefaultGpioValue()))
		if err != nil {
			log.Fatal(err)
		}
		HeaterLine = line
		return line
	}
	return HeaterLine
}
