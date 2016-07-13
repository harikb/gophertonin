package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"

	"io/ioutil"
	"log"
	"net/http"
)

var button *gpio.GroveButtonDriver
var blue *gpio.GroveLedDriver
var green *gpio.GroveLedDriver
var red *gpio.GroveLedDriver
var rotary *gpio.GroveRotaryDriver
var sensor *gpio.GroveTemperatureSensorDriver
var light *gpio.GroveLightSensorDriver

var currentLightState = 2700

func DetectLight(level int) {
	if level >= 400 {
		fmt.Println("Light detected")
		TurnOff()
		blue.On()
		<-time.After(1 * time.Second)
		Reset()
	}
}

func TurnOff() {
	blue.Off()
	green.Off()
}

func Reset() {
	TurnOff()
	fmt.Println("Airlock ready.")
	green.On()
}

func ToggleScreenTemp() {
	if currentLightState == 2700 {
		currentLightState = 6500
	} else {
		currentLightState = 2700
	}
	resp, err := http.Get(fmt.Sprintf("http://10.56.241.26:8001/gophertonin?light=%d", currentLightState))
	if err != nil {
		log.Printf("Error sending request to gophertonin server")
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("%s", body)
}

func main() {
	gbot := gobot.NewGobot()

	a := api.NewAPI(gbot)
	a.Start()

	// digital
	board := edison.NewEdisonAdaptor("edison")

	button = gpio.NewGroveButtonDriver(board, "button", "2")
	blue = gpio.NewGroveLedDriver(board, "blue", "3")
	green = gpio.NewGroveLedDriver(board, "green", "4")
	red = gpio.NewGroveLedDriver(board, "red", "5")

	// analog
	rotary = gpio.NewGroveRotaryDriver(board, "rotary", "0")
	sensor = gpio.NewGroveTemperatureSensorDriver(board, "sensor", "1")
	light = gpio.NewGroveLightSensorDriver(board, "light", "3")

	work := func() {
		Reset()

		gobot.On(button.Event(gpio.Push), func(data interface{}) {
			log.Printf("Button press")
			ToggleScreenTemp()
			log.Printf("Action for button press done")
		})

		gobot.On(button.Event(gpio.Release), func(data interface{}) {
			Reset()
		})

		gobot.On(rotary.Event("data"), func(data interface{}) {
			fmt.Println("rotary", data)
		})

		gobot.On(light.Event("data"), func(data interface{}) {
			DetectLight(data.(int))
		})

	}

	robot := gobot.NewRobot("airlock",
		[]gobot.Connection{board},
		[]gobot.Device{button, blue, green, red, rotary, sensor, light},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
