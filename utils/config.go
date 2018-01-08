package utils

import (
	"github.com/larspensjo/config"
	"os"
	//"fmt"
)

func LoadConfig(configName string, innerConfigName string) (result map[string]string, err error) {

	var TOPIC = make(map[string]string)

	configFile, err := os.Getwd()
	if err != nil {
		return TOPIC, err
	}
	cfg, err := config.ReadDefault(configFile + "/config/" + configName + ".ini")
	if err != nil {
		return TOPIC, err
	}
	if cfg.HasSection(innerConfigName) {
		section, err := cfg.SectionOptions(innerConfigName)
		if err == nil {
			for _, v := range section {
				options, err := cfg.String(innerConfigName, v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
		return TOPIC, nil
	}
	return TOPIC, err
}
