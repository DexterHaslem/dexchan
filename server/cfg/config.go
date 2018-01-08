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

package cfg

import (
	"encoding/json"
	"io/ioutil"
)

// C is the configuration options for server to use
type C struct {
	StaticDir  string `json:"staticDir"`
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
