package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

type lightHandlder struct{}

type response struct {
	Error  string `json:"message,omitempty"`
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/gophertonin", handleRequest)
	http.ListenAndServe("10.56.241.26:8001", nil)
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
