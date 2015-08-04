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
	//"github.com/jordansissel/fingerpoken/mdp"
	//czmq "github.com/jordansissel/goczmq"
	log "github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"
	"github.com/jordansissel/fingerpoken/mdp"
	"github.com/jordansissel/fingerpoken/workers"
	czmq "github.com/jordansissel/goczmq"
	"os"
)

type Settings struct {
	workers.Settings
}

const service string = "oauth"

func main() {
	// Flags:
	//   --broker host:port
	//   --broker-public-key abcdef...

	var settings Settings
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

	// Load the client cert
	cert, err := czmq.NewCertFromFile(settings.ClientCertificatePath)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// Load the server cert
	server_cert, err := czmq.NewCertFromFile(settings.BrokerKeyPath)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.WithFields(log.Fields{"broker": settings.Broker, "service": service}).Info("Worker starting up")
	worker := mdp.NewJSONRPCWorker(settings.Broker, service)
	worker.CurveServerPublicKey = server_cert.PublicText()
	worker.CurveCertificate = cert
	worker.Register(&OAuth{})
	worker.Run()
}
