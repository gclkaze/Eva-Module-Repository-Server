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

func InitWithPropertiesPath(prop string) {
	if TheConfigReader != nil {
		return
	}
	TheConfigReader = &ConfigReader{}
	TheConfigReader.error = TheConfigReader.readPropertiesFromPath(prop)
	if TheConfigReader.error != nil {
		TheConfigReader.onError = true
	}
}

func InitWithPropertiesMap(m *map[string]string) {
	if TheConfigReader != nil {
		return
	}
	TheConfigReader = &ConfigReader{}
	TheConfigReader.error = TheConfigReader.readPropertiesFromMap(m)
	if TheConfigReader.error != nil {
		TheConfigReader.onError = true
	}
}

func (c *ConfigReader) readPropertiesFromPath(prop string) error {
	c.properties = properties.MustLoadFile(prop, properties.UTF8)
	if c.properties == nil {
		return fmt.Errorf("couldn't read application properties file %s", prop)
	}
	return nil
}

func (c *ConfigReader) readPropertiesFromMap(m *map[string]string) error {
	c.properties = properties.LoadMap(*m)
	if c.properties == nil {
		return fmt.Errorf("couldn't read application properties map")
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
