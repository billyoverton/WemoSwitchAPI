package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/billyoverton/go.wemo"
	"github.com/gorilla/mux"
)

// WemoPort default port used by WeMo Devices
const WemoPort = 49153

var switches Switches

func main() {
	configPath := flag.String("config", "./switches.json", "path to configuration file")
	port := flag.Int("port", 8080, "Port to run the server on")

	flag.Parse()

	configuration, err := parseConfig(*configPath)

	if err != nil {
		fmt.Println("Error opening configuration file: ", err)
		os.Exit(1)
	}

	for _, configSwitch := range configuration.Switches {
		switches = append(switches, Switch{DeviceName: configSwitch.Name, Device: wemo.Device{Host: fmt.Sprintf("%v:%v", configSwitch.IP, WemoPort)}})
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/switches", SwitchesIndex).Methods("GET")
	router.HandleFunc("/switches/{name}", SwitchGet).Methods("GET")
	router.HandleFunc("/switches/{name}/{command}", SwitchCommand).Methods("POST")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), router))
}

// Index handler
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// SwitchesIndex index for all switches
func SwitchesIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var switchsReturn SwitchesReturn

	for _, switchDevice := range switches {
		switchsReturn.Switches = append(switchsReturn.Switches, GetSwitchReturn(switchDevice))
	}

	if err := json.NewEncoder(w).Encode(switchsReturn); err != nil {
		panic(err)
	}
}

// SwitchGet gets details on a specific switch
func SwitchGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switchName := vars["name"]

	var foundSwitch Switch
	found := false

	for _, switchDevice := range switches {
		if switchName == switchDevice.DeviceName {
			foundSwitch = switchDevice
			found = true
			break
		}
	}

	if !found {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "{\"error\": \"Switch not found\"}")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(GetSwitchReturn(foundSwitch)); err != nil {
		panic(err)
	}
}

// SwitchCommand control commands for a switch
func SwitchCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switchName := vars["name"]
	command := vars["command"]

	var foundSwitch Switch
	found := false

	for _, switchDevice := range switches {
		if switchName == switchDevice.DeviceName {
			foundSwitch = switchDevice
			found = true
			break
		}
	}

	if !found {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "{\"error\": \"Switch not found\"}")
		return
	}

	switch command {
	case "on":
		foundSwitch.Device.On()
	case "off":
		foundSwitch.Device.Off()
	case "toggle":
		foundSwitch.Device.Toggle()
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "{\"error\": \"Command not found\"}")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(GetSwitchReturn(foundSwitch)); err != nil {
		panic(err)
	}

}
