package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Bridge struct {
	XMLName      xml.Name `xml:"root"`
	BaseUrl      string   `xml:"URLBase"`
	DeviceName   string   `xml:"device>modelName"`
	DeviceModel  string   `xml:"device>modelNumber"`
	SerialNumber string   `xml:"device>serialNumber"`
}

type RegisterUser struct {
	//Username   string `json:"username"`
	DeviceName string `json:"devicetype"`
}

type Success struct {
	Username string `json:"username,omitempty"`
}

type Error struct {
	ErrorType        int    `json:"type,omitempty"`
	ErrorAddress     string `json:"address,omitempty"`
	ErrorDescription string `json:"description,omitempty"`
}

type Response struct {
	Success Success `json:"success,omitempty"`
	Error   Error   `json:"error,omitempty"`
}

type ApiResponse []Response

type Light struct {
	Config struct {
		Type string `json:"archetype"`
	} `json:"config"`
	Name string `json:"name"`
}

type Lights map[string]Light

type LightState struct {
	On               bool      `json:"on"`
	Brightness       uint8     `json:"bri,omitempty"`
	Hue              uint16    `json:"hue,omitempty"`
	Saturation       uint8     `json:"sat,omitempty"`
	Coordinates      []float64 `json:"xy,omitempty"`
	ColorTemperature uint16    `json:"ct,omitempty"`
	TransitionTime   uint16    `json:"transitiontime,omitempty"`
}

func (resp ApiResponse) HasErrors() bool {
	hasErrors := false

	for _, r := range resp {
		if r.Error.ErrorDescription != "" {
			hasErrors = true
			break
		}
	}

	return hasErrors
}

func (resp ApiResponse) PrintErrors() {
	for _, r := range resp {
		if r.Error.ErrorDescription != "" {
			log.Printf("hue: bridge: %s\n", r.Error.ErrorDescription)
		}
	}
}

// RegisterUser registers a new user (username randomly generated as per hue
// api definition) and returns the username it is to be used for creating an
// developer interface on the hue bridge; and generally should be called as
// one of the first things in the execution (for any API request)
func (b Bridge) RegisterUser(deviceName string) string {
	body, err := json.Marshal(RegisterUser{deviceName})
	if err != nil {
		logError(err.Error())
	}

	resp, err := http.Post(b.BaseUrl+"api", "application/json", bytes.NewBuffer(body))
	if err != nil {
		logError(err.Error())
	}
	defer resp.Body.Close()

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logError(readErr.Error())
	}

	var apiResp ApiResponse
	err = json.Unmarshal(respBody, &apiResp)
	if err != nil {
		logError(err.Error())
	}

	if !apiResp.HasErrors() {
		return apiResp[0].Success.Username
	} else {
		apiResp.PrintErrors()
	}

	return ""
}

func (b Bridge) AllLights(username string) Lights {
	resp, err := http.Get(b.BaseUrl + "api/" + username + "/lights")
	if err != nil {
		logError(err.Error())
	}
	defer resp.Body.Close()

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logError(readErr.Error())
	}

	var lights Lights
	err = json.Unmarshal(respBody, &lights)
	if err != nil {
		logError(err.Error())
	}

	return lights
}

func (b Bridge) SetLightState(username string, id string, state LightState) {
	body, bodyErr := json.Marshal(state)
	if bodyErr != nil {
		logError(bodyErr.Error())
	}

	fmt.Printf("hue: set light '%s'\n", id)
	req, reqErr := http.NewRequest(http.MethodPut, b.BaseUrl+"api/"+username+"/lights/"+id+"/state", bytes.NewBuffer(body))
	if reqErr != nil {
		logError(reqErr.Error())
	}

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		logError(respErr.Error())
	}
	defer resp.Body.Close()

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logError(readErr.Error())
	}

	var apiResp ApiResponse
	err := json.Unmarshal(respBody, &apiResp)
	if err != nil {
		logError(err.Error())
	}

	if apiResp.HasErrors() {
		apiResp.PrintErrors()
	}
}
