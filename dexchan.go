// dexChan copyright Dexter Haslem <dmh@fastmail.com> 2018
// This file is part of dexChan
//
// dexChan is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// dexChan is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with dexChan. If not, see <http://www.gnu.org/licenses/>.


package main

import (
	"dexchan/server/cfg"
	"dexchan/server"
	"log"
	"dexchan/server/db"
)

func loadCfg() *cfg.C {
	//  todo: flag for passing in config
	c, err := cfg.From("dexchan.json")
	if err != nil {
		return &cfg.C{
			StaticDir:  "static",
			Port:       8080,
			DbUsername: "dexchan",
			DbPassword: "dexchan",
			DbHost:     "localhost",
			DbName:     "dexchan",
			DbPort:     5432,
		}
	}
	return c
}

func main() {
	c := loadCfg()
	s := server.Server{Config: c}

	d, err := db.Open(s.Config)
	if err != nil {
		log.Fatalf("failed to connect to database: %s\n", err)
	}
	s.Db = d

	log.Printf("server starting. config=%v\n", s.Config)
	s.Start()
}
