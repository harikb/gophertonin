package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	dprofile "github.com/pkg/profile"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var appVersion string

type cmdArgs struct {
	cpuProfile    bool
	memProfile    bool
	verbose       bool
	serverAddress string
}

var button *gpio.GroveButtonDriver

var currentLightState = 2700

func ToggleScreenTemp(gophertoninServer string) {
	if currentLightState == 2700 {
		currentLightState = 6500
	} else {
		currentLightState = 2700
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/gophertonin?light=%d", gophertoninServer, currentLightState))
	if err != nil {
		log.Printf("Error sending request to gophertonin server")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("%s", body)
}

func processCmdline() cmdArgs {

	var args cmdArgs

	// usage is customized to include Version number
	var usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (Version: %s):\n", os.Args[0], appVersion)
		flag.PrintDefaults()
	}

	flag.BoolVarP(&args.cpuProfile, "cpu-profile", "", false,
		"(for debug only) CPU profile this run")
	flag.BoolVarP(&args.memProfile, "mem-profile", "", false,
		"(for debug only) MEM profile this run")
	flag.StringVarP(&args.serverAddress, "", "s", "0.0.0.0:8001",
		"(for debug only) MEM profile this run")
	flag.Usage = usage
	flag.Parse()

	return args
}

func main() {

	args := processCmdline()

	if args.cpuProfile {
		defer dprofile.Start(dprofile.CPUProfile).Stop()
	} else if args.memProfile {
		defer dprofile.Start(dprofile.MemProfile).Stop()
	}

	gbot := gobot.NewGobot()

	a := api.NewAPI(gbot)
	a.Start()

	// digital
	board := edison.NewEdisonAdaptor("edison")

	button = gpio.NewGroveButtonDriver(board, "button", "2")

	work := func() {

		gobot.On(button.Event(gpio.Push), func(data interface{}) {
			log.Printf("Button press")
			ToggleScreenTemp(args.serverAddress)
			log.Printf("Action for button press done")
		})

	}

	robot := gobot.NewRobot("airlock",
		[]gobot.Connection{board},
		[]gobot.Device{button},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
