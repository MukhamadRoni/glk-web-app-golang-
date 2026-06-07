package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Menu struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Icon     string `json:"icon,omitempty"`
	Children []Menu `json:"children,omitempty"`
}

var AdminMenus []Menu

func LoadMenu(filePath string) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("[WARN] Could not read menu file: %v\n", err)
		return
	}
	if err := json.Unmarshal(bytes, &AdminMenus); err != nil {
		log.Printf("[WARN] Could not parse menu file: %v\n", err)
	}
}
