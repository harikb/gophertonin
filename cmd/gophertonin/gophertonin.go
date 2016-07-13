package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	dprofile "github.com/pkg/profile"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var appVersion string

type cmdArgs struct {
	cpuProfile  bool
	memProfile  bool
	verbose     bool
	bindAddress string
}

type response struct {
	Error  string `json:"message,omitempty"`
	Status string `json:"status"`
}

func setFluxValue(lv int64) error {

	cmd := exec.Command("defaults", "write", "org.herf.Flux", "dayColorTemp", "-int", fmt.Sprintf("%d", lv))
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Can't invoke cmd: %v", err)
	}
	if stdout.String() != "" || stderr.String() != "" {
		return fmt.Errorf("%s, %s", stdout.String(), stderr.String())
	}
	return err
}

func restartFlux() error {

	cmd := exec.Command("killall", "Flux")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Can't invoke cmd: %v", err)
	}
	if stdout.String() != "" || stderr.String() != "" {
		return fmt.Errorf("%s, %s", stdout.String(), stderr.String())
	}

	cmd = exec.Command("open", "/Applications/Flux.app")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Can't invoke cmd: %v", err)
	}
	if stdout.String() != "" || stderr.String() != "" {
		return fmt.Errorf("%s, %s", stdout.String(), stderr.String())
	}
	return err
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	log.Print(spew.Sdump(r))

	lvs := r.FormValue("light")
	log.Printf("Value for lvs %s", lvs)
	lv, err := strconv.ParseInt(lvs, 10, 32)
	if err != nil || lv < 2700 || lv > 6500 {
		log.Printf("Incorrect value for lv %d", lv)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := response{}
	resp.Status = "Success"

	err = setFluxValue(lv)
	if err != nil {
		resp.Error = fmt.Sprintf("%s", err)
		resp.Status = "Fail"
	} else {
		err = restartFlux()
		if err != nil {
			resp.Error = fmt.Sprintf("%s", err)
			resp.Status = "Fail"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	respBytes, err := json.Marshal(resp)
	if err != nil {
		respBytes = []byte(fmt.Sprintf("{\"status\": \"Fail\", \"error\": \"Marshal error: %v\"}", err))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
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
	flag.StringVarP(&args.bindAddress, "", "b", "0.0.0.0:8001",
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

	http.HandleFunc("/gophertonin", handleRequest)
	log.Printf("Listening to %s", args.bindAddress)
	log.Fatal(http.ListenAndServe(args.bindAddress, nil))
}
