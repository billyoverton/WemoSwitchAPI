package main

import "github.com/billyoverton/go.wemo"

// Switch Wemo device
type Switch struct {
	DeviceName string
	Device     wemo.Device
}

// Switches a collection of switches
type Switches []Switch

// SwitchReturn json representation for a switch
type SwitchReturn struct {
	DeviceName string `json:"name"`
	State      string `json:"state"`
}

// SwitchesReturn json representation for multiple switches
type SwitchesReturn struct {
	Switches []SwitchReturn `json:"switches"`
}

// Get the json return for a given switch
func GetSwitchReturn(switchDevice Switch) SwitchReturn {
	var switchReturn SwitchReturn
	switchReturn.DeviceName = switchDevice.DeviceName

	switch switchDevice.Device.GetBinaryState() {
	case -1:
		switchReturn.State = "Error Getting State"
	case 0:
		switchReturn.State = "off"
	case 1:
		switchReturn.State = "on"
	}

	return switchReturn

}
