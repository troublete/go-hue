package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/aarzilli/golua/lua"
	"github.com/koron/go-ssdp"
)

func logError(err string) {
	log.Fatalf("hue: %s", err)
}

func main() {
	interval := flag.Int64("interval", -1, "Set the interval in milliseconds in which the provided script is run.")
	flag.Parse()

	list, err := ssdp.Search(ssdp.All, 1, "")
	if err != nil {
		logError(err.Error())
	}

	regex := regexp.MustCompile(".*\\:80\\/description\\.xml$")
	locations := make(map[string]Bridge)

	regexBridge := regexp.MustCompile("hue bridge")

	for _, srv := range list {
		if regex.MatchString(srv.Location) {
			resp, respErr := http.Get(srv.Location)
			if respErr != nil {
				logError(respErr.Error())
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logError(err.Error())
			}

			b := Bridge{}
			err = xml.Unmarshal(body, &b)
			if err != nil {
				logError(err.Error())
			}

			if regexBridge.MatchString(b.DeviceName) && b.DeviceModel == "BSB002" {
				locations[b.SerialNumber] = b
			}
		}
	}

	script := os.Args[1]
	if script == "" {
		logError("script must be provided to be executed.")
	}

	l := lua.NewState()
	l.OpenLibs()

	l.NewTable()
	for serial, b := range locations {
		l.PushString(serial)
		l.NewTable()
		l.PushString("base_url")
		l.PushString(b.BaseUrl)
		l.SetTable(-3)
		l.PushString("register_user")
		l.PushGoFunction(func(l *lua.State) int {
			device := l.CheckString(1)
			username := b.RegisterUser(device)
			l.PushString(username)
			return 1
		})
		l.SetTable(-3)
		l.PushString("lights")
		l.PushGoFunction(func(l *lua.State) int {
			username := l.CheckString(1)

			l.NewTable()
			for i, light := range b.AllLights(username) {
				id := i
				l.PushString(id)
				l.NewTable()
				l.PushString("type")
				l.PushString(light.Config.Type)
				l.SetTable(-3)
				l.PushString("name")
				l.PushString(light.Name)
				l.SetTable(-3)
				l.PushString("turn_off")
				l.PushGoFunction(func(l *lua.State) int {
					state := &LightState{}
					state.On = false
					b.SetLightState(username, id, *state)
					return 0
				})
				l.SetTable(-3)
				l.PushString("turn_on")
				l.PushGoFunction(func(l *lua.State) int {
					state := &LightState{}
					state.On = true
					b.SetLightState(username, id, *state)
					return 0
				})
				l.SetTable(-3)
				l.PushString("set")
				l.PushGoFunction(func(l *lua.State) int {
					hue := l.CheckInteger(1)
					bri := l.CheckInteger(2)
					sat := l.CheckInteger(3)
					temp := l.CheckInteger(4)
					tt := l.CheckInteger(5)

					state := &LightState{}
					state.On = true

					if hue > -1 {
						state.Hue = uint16(hue)
					}

					if bri > -1 {
						state.Brightness = uint8(bri)
					}

					if sat > -1 {
						state.Saturation = uint8(sat)
					}

					if temp > -1 {
						state.ColorTemperature = uint16(temp)
					}

					if tt > -1 {
						state.TransitionTime = uint16(tt)
					}

					b.SetLightState(username, id, *state)
					return 0
				})
				l.SetTable(-3)
				l.SetTable(-3)
			}
			return 1
		})
		l.SetTable(-3)
		l.SetTable(-3)
	}
	l.SetGlobal("bridges")

	if (*interval) == -1 {
		if err := l.DoFile(script); err != nil {
			logError(err.Error())
		}
	} else {
		c := 1
		for {
			l.PushInteger(int64(c))
			l.SetGlobal("loop")

			if err := l.DoFile(script); err != nil {
				logError(err.Error())
			}

			time.Sleep(time.Millisecond * time.Duration(*interval))
			c = c + 1
		}
	}
}
