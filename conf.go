package mdoc

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/zrcoder/mdoc/internal/log"
	"github.com/zrcoder/mdoc/internal/model"
)

var cfg *model.Config

// InitWithFile initializes configuration from conf assets and given custom
// configuration file. If `customConf` is empty, it falls back to default
// location, i.e. "<WORK DIR>/custom".
func InitWithFile(customConfigFile string) (err error) {
	log.Info("custom config file:", customConfigFile)
	if customConfigFile == "" {
		return nil
	}

	data, err := os.ReadFile(customConfigFile)
	if err != nil {
		return err
	}
	current := GetConfig()
	err = yaml.Unmarshal(data, current)
	if err != nil {
		return err
	}
	current.DocsBasePath = strings.TrimRight(current.DocsBasePath, "/")
	log.Info(current.DocsBasePath)
	if current.DocsBasePath != "" && current.DocsBasePath[0] != '/' {
		return errors.New("invalid docsBasePath, should start with '/'")
	}
	cfg = current
	return nil
}

func GetConfig() *model.Config {
	if cfg != nil {
		return cfg
	}
	cfg = &model.Config{
		HttpAddr:      "localhost",
		HttpPort:      "9999",
		DocsDirectory: ".",
	}
	return cfg
}
