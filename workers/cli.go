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
package workers

import (
	log "github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"
	"os"
)

type Settings struct {
	Broker                string `long:"broker" required:"true"`
	BrokerKeyPath         string `long:"broker-public-key" required:"true"`
	ClientCertificatePath string `long:"client-certificate" required:"true"`
}

func ParseArgs(settings interface{}) []string {
	remaining, err := flags.Parse(settings)
	if err != nil {
		switch err.(*flags.Error).Type {
		case flags.ErrHelp:
			os.Exit(0)
		default:
			log.Error(err.Error())
			os.Exit(1)
		}
	}

	return remaining
}
