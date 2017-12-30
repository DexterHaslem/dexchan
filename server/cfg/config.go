package cfg

import (
	"encoding/json"
	"io/ioutil"
)

// C is the configuration options for server to use
type C struct {
	DbName     string `json:"dbName"`
	DbHost     string `json:"dbHost"`
	DbPort     int    `json:"dbPort"`
	DbUsername string `json:"dbUsername"`
	DbPassword string `json:"dbPassword"`
	Port       int    `json:"port"`
}

// From will try to load a config object from a file in json format
func From(file string) (*C, error) {
	fc, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ret := &C{}
	err = json.Unmarshal(fc, ret)
	return ret, err
}

// Save will save the current configuration object to a given filename in
// json format
func (c *C) Save(file string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, b, 0666)
}
