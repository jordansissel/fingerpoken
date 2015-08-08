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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

func init() {
	gomniauth.SetSecurityKey(signature.RandomKey(64))
}

type OAuth struct {
	ClientKeys map[string]ClientKey // key == provider name
}

type ClientKey struct {
	ClientId     string `json:id`
	ClientSecret string `json:secret`
}

type OAuthLoginStart struct {
	Provider    string `json:provider`
	Scope       string `json:scope`
	CallbackURL string `json:callback`
}

type OAuthLoginRedirect struct {
	AuthURL string `json:authurl`
}

type Empty struct{}
type ProviderList []string

// RPC Call OAuth.ListProviders
func (o *OAuth) ListProviders(args *Empty, reply *ProviderList) error {
	// TODO(sissel): Make this actually look up the known providers list.
	*reply = append(*reply, "google")
	return nil
}

// RPC Call OAuth.Login
// Takes an OAuthLoginStart w/ provider name
func (o *OAuth) StartLogin(login *OAuthLoginStart, reply *OAuthLoginRedirect) (err error) {
	provider := google.New("952973200910-mt6fvajdvnhjgb9h7hlli2sqpjmu4lp5.apps.googleusercontent.com", "6J2oSU53bKuSmzQIpC8sEH3f", login.CallbackURL)
	state := gomniauth.NewState("after", "success")
	options := objx.Map{
		"scope": login.Scope,
	}
	reply.AuthURL, err = provider.GetBeginAuthURL(state, options)
	return nil
}
