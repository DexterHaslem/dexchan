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
