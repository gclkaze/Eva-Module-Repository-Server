package config

import (
	"fmt"
	"os"

	"github.com/magiconair/properties"
)

type ConfigReader struct {
	properties *properties.Properties
	onError    bool
	error      error
}

var TheConfigReader *ConfigReader = nil

func Init() {
	if TheConfigReader != nil {
		return
	}
	TheConfigReader = &ConfigReader{}
	TheConfigReader.error = TheConfigReader.readProperties()
	if TheConfigReader.error != nil {
		TheConfigReader.onError = true
	}
}

func InitWithProperties(prop string) {
	if TheConfigReader != nil {
		return
	}
	TheConfigReader = &ConfigReader{}
	TheConfigReader.error = TheConfigReader.readPropertiesFromString(prop)
	if TheConfigReader.error != nil {
		TheConfigReader.onError = true
	}
}

func (c *ConfigReader) readPropertiesFromString(prop string) error {
	c.properties = properties.MustLoadFile(prop, properties.UTF8)
	if c.properties == nil {
		return fmt.Errorf("couldn't read application properties file %s", prop)
	}
	return nil
}

func (c *ConfigReader) readProperties() error {
	currentPath, _ := os.Getwd()
	f := currentPath + "\\internal\\config\\application.properties"
	c.properties = properties.MustLoadFile(f, properties.UTF8)
	if c.properties == nil {
		return fmt.Errorf("couldn't read application properties file %s", f)
	}
	return nil
}

func (c ConfigReader) IsOnError() bool {
	return c.onError
}

func (c *ConfigReader) GetProperties() *properties.Properties {
	return c.properties
}
