package Config

import (
	"encoding/json"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"os"
)

// just loads configValues into Config struct
func LoadConfigValues(file string) (*Models.Config, error) {
	configFile, err := os.Open(file)

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			return
		}
	}(configFile)

	if err != nil {
		return nil, err
	}

	var config Models.Config

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
