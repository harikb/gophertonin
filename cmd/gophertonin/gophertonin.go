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

func handleRequest(w http.ResponseWriter, r *http.Request) {

	log.Print(spew.Sdump(r))

	lvs := r.FormValue("light")
	log.Printf("Value for lvs %s", lvs)
	lv, err := strconv.ParseInt(lvs, 10, 32)
	if err != nil {
		log.Printf("Incorrect value for lv %d", lv)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := response{}
	resp.Status = "Success"

	cmd := exec.Command("./fluxcontrol.bash", fmt.Sprintf("%d", lv))
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		resp.Error = fmt.Sprintf("Can't invoke cmd: %v", err)
		resp.Status = "Fail"
	}
	if stdout.String() != "" || stderr.String() != "" {
		resp.Error = fmt.Sprintf("%s, %s", stdout.String(), stderr.String())
		resp.Status = "Fail"
	}

	w.Header().Set("Content-Type", "application/json")
	respBytes, err := json.Marshal(resp)
	if err != nil {
		respBytes = []byte(fmt.Sprintf("{\"status\": \"Fail\", \"error\": \"Marshal error: %v\"}", err))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
