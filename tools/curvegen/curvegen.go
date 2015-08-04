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
	log "github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"
	czmq "github.com/jordansissel/goczmq"
	"os"
)

type Settings struct {
	OutputPath string `short:"o"`
}

func main() {
	cert := czmq.NewCert()
	defer cert.Destroy()

	var settings Settings
	settings.OutputPath = "./curve"
	_, err := flags.Parse(&settings)
	if err != nil {
		switch err.(*flags.Error).Type {
		case flags.ErrHelp:
			os.Exit(0)
		default:
			log.Error(err.Error())
			os.Exit(1)
		}
	}

	cert.Save(settings.OutputPath)
	log.WithFields(log.Fields{"path": settings.OutputPath}).Info("Public cert saved")
	log.WithFields(log.Fields{"path": fmt.Sprintf("%s_secret", settings.OutputPath)}).Info("Secret cert saved")
}
