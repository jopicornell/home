package main

import (
	"bufio"
	"github.com/jopicornell/thermonick/log"
	"github.com/jopicornell/thermonick/models"
	"github.com/jopicornell/thermonick/utils"
	"github.com/warthog618/gpiod"
	"github.com/yryz/ds18b20"
	"os"
	"time"
)

type Status struct {
	HeaterOn bool
	LightOn  bool

	BaskingTemp float64
	ColdTemp    float64
}

var sensorNameMap = map[string]string{
	"basking": "28-012032eb3734",
	"cold":    "28-012032cd28ce",
}

var currentStatus = Status{
	HeaterOn:    false,
	LightOn:     false,
	BaskingTemp: 0,
	ColdTemp:    0,
}

var HeaterLine *gpiod.Line
var HeaterLineNumber = 14
var LightLine *gpiod.Line
var LightLineNumber = 15

func main() {
	log.CreateLoggers()
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}
	log.Logger.Infof("sensor IDs: %v", sensors)
	ReadCommandsRoutine()
	for {
		baskingTemperature, err := getBaskingTemp()
		if err != nil {
			log.Logger.Errorf("Error getting basking temperature: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		coldTemperature, err := getColdTemp()
		if err != nil {
			log.Logger.Errorf("Error getting basking temperature: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		ChangeStatusTemperatures(baskingTemperature, coldTemperature)
		if ShouldActivateHeater(baskingTemperature) {
			ActivateHeater(baskingTemperature)
		} else if ShouldDeactivateHeater(baskingTemperature) {
			DeactivateHeater(baskingTemperature)
		}

		if ShouldActivateLight() {
			ActivateLight()
		} else if ShouldDeactivateLight() {
			DeactivateLight()
		}

		time.Sleep(time.Minute)
	}
	//	response, err := http.Post("https://umczz0pvpc.execute-api.eu-central-1.amazonaws.com/prod/temperatures", "application/json", bytes.NewBuffer(body))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("Created temperatures. Response: %+v", response)
}

func ChangeStatusTemperatures(baskingTemperature float64, coldTemperature float64) {
	if baskingTemperature != currentStatus.BaskingTemp {
		log.Logger.Debugf("Changed baskingTemperature: %.2f", baskingTemperature)
		currentStatus.BaskingTemp = baskingTemperature
	}
	if coldTemperature != currentStatus.ColdTemp {
		log.Logger.Debugf("Changed cold temperature: %.2f", coldTemperature)
		currentStatus.ColdTemp = coldTemperature
	}
}

func ReadCommandsRoutine() {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			cmd, err := reader.ReadString('\n')
			if err != nil {
				log.Logger.Fatal(err)
			}
			sunrise, sunset := utils.AdjustedSunriseSunset(39.57, 2.65, time.Now(), getCurrentCondition().GetMinHours(), getCurrentCondition().GetMaxHours())
			if cmd == "status\n" {
				log.Logger.Debugf("Current status \n Heater: %v \n Light: %v \n BaskingTemp: %f \n ColdTemp: %f \n Sunrise: %s\n Sunset: %s \n Light hours %f \n", currentStatus.HeaterOn, currentStatus.LightOn, currentStatus.BaskingTemp, currentStatus.ColdTemp, sunrise, sunset, sunset.Sub(sunrise).Hours())
			}
			if cmd == "heater on\n" {
				ActivateHeater(currentStatus.BaskingTemp)
			}
			if cmd == "heater off\n" {
				DeactivateHeater(currentStatus.BaskingTemp)
			}
			if cmd == "light on\n" {
				ActivateLight()
			}
			if cmd == "light off\n" {
				DeactivateLight()
			}
			if cmd == "exit\n" {
				os.Exit(0)
			}
		}

	}()
}

func DeactivateHeater(baskingTemperature float64) {
	log.Logger.Infof("Deactivating heater temperature reached %.2f more than %.2f", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getHeaterLine().SetValue(0)
	if err != nil {
		log.Logger.Errorf("Error deactivating heater: %s", err.Error())
	}
	currentStatus.HeaterOn = false
}

func ActivateHeater(baskingTemperature float64) {
	log.Logger.Infof("Activating heater temperature reached %.2f less than %.2f", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getHeaterLine().SetValue(1)
	if err != nil {
		log.Logger.Errorf("Error activating heater: %+v", err)
	}
	currentStatus.HeaterOn = true
}

func DeactivateLight() {
	log.Logger.Printf("Deactivating light due to sunset")
	err := getLightLine().SetValue(0)
	if err != nil {
		log.Logger.Errorf("Error deactivating light: %s", err.Error())
	}
	currentStatus.LightOn = false
}

func ActivateLight() {
	log.ChangesLogger.Printf("Activating light due to sunrise")
	err := getLightLine().SetValue(1)
	if err != nil {
		log.Logger.Errorf("Error activating light: %+v", err)
	}
	currentStatus.LightOn = true
}

func getBaskingTemp() (float64, error) {
	return getTemperature(sensorNameMap["basking"])
}

func getColdTemp() (float64, error) {
	return getTemperature(sensorNameMap["cold"])
}

func getTemperature(sensorName string) (float64, error) {
	return ds18b20.Temperature(sensorName)
}

func ShouldActivateHeater(temperature float64) bool {
	condition := getCurrentCondition()
	return !currentStatus.HeaterOn && condition.IsHeaterOn(time.Now(), temperature)
}

func ShouldDeactivateHeater(temperature float64) bool {
	condition := getCurrentCondition()
	return currentStatus.HeaterOn && condition.IsHeaterOff(time.Now(), temperature)
}

func ShouldActivateLight() bool {
	condition := getCurrentCondition()
	return !currentStatus.LightOn && condition.IsLightOn(time.Now())
}

func ShouldDeactivateLight() bool {
	condition := getCurrentCondition()
	return currentStatus.LightOn && !condition.IsLightOn(time.Now())
}

func getCurrentCondition() models.Condition {
	return &models.SunsetSunriseCondition{
		MinimumHours:     8,
		MaximumHours:     12,
		TemperatureNight: 21,
		TemperatureDay:   31,
		Latitude:         39.57,
		Longitude:        2.65,
	}
}

func getDefaultGpioValue() int {
	condition := getCurrentCondition()
	baskingTemperature, err := getBaskingTemp()
	if err != nil {
		log.Logger.Errorf("Error getting basking temperature: %s", err.Error())
		return 0
	}
	if condition.IsHeaterOn(time.Now(), baskingTemperature) {
		return 1
	}
	return 0
}

func getHeaterLine() *gpiod.Line {
	if HeaterLine == nil {
		line, err := gpiod.RequestLine("gpiochip0", HeaterLineNumber, gpiod.AsOutput(getDefaultGpioValue()))
		if err != nil {
			log.Logger.Fatal(err)
		}
		HeaterLine = line
		return line
	}
	return HeaterLine
}

func getLightLine() *gpiod.Line {
	if LightLine == nil {
		line, err := gpiod.RequestLine("gpiochip0", LightLineNumber, gpiod.AsOutput(getDefaultGpioValue()))
		if err != nil {
			log.Logger.Fatal(err)
		}
		LightLine = line
		return line
	}
	return HeaterLine
}
