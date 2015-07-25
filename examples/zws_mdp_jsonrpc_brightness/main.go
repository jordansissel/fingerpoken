// This file is part of fingerpoken
// Copyright (C) 2015 Jordan Sissel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"fmt"
	"github.com/jordansissel/fingerpoken/mdp"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Screen struct{}

type Empty struct{}
type BrightnessSetting struct {
	Percent float64 `json:percent`
}

// TODO(sissel): Read this at runtime
const MAX_BRIGHTNESS = 1388

func (s *Screen) SetBrightness(args *BrightnessSetting, reply *Empty) error {
	log.Printf("Setting brightness to %f", args.Percent)

	brightness := int(args.Percent * MAX_BRIGHTNESS)
	data := []byte(fmt.Sprintf("%d", brightness))
	//err := ioutil.WriteFile("/sys/class/backlight/intel_backlight/brightness", data, 0600)
	fd, err := os.OpenFile("/sys/class/backlight/intel_backlight/brightness", os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write(data)
	// TODO(sissel): Check return bytes read.
	if err != nil {
		return err
	}
	*reply = Empty{}
	return nil
}

func (s *Screen) GetBrightness(args *Empty, reply *BrightnessSetting) error {
	data, err := ioutil.ReadFile("/sys/class/backlight/intel_backlight/brightness")
	if err != nil {
		return err
	}
	value, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		return err
	}

	reply.Percent = float64(value) / MAX_BRIGHTNESS
	return nil
}

func main() {
	broker_endpoint := os.Args[1]
	service := os.Args[2]
	w := mdp.NewJSONRPCWorker(broker_endpoint, service)
	w.Register(&Screen{})
	w.Run()
}
