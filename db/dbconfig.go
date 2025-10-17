package db

import (
	"fmt"
	"os"

	"github.com/magiconair/properties"
)

type EvaModuleRepositoryDatabaseConfig struct {
	username         string
	pwd              string
	host             string
	port             int
	dbName           string
	connectionString string
}

func NewEvaModuleRepositoryDatabaseConfig() *EvaModuleRepositoryDatabaseConfig {
	inst := &EvaModuleRepositoryDatabaseConfig{}
	return inst
}

func (c *EvaModuleRepositoryDatabaseConfig) LoadFromProperties(p *properties.Properties) error {
	c.port = p.GetInt("db_port", 0)
	c.host = p.GetString("db_hostname", "")
	c.dbName = p.GetString("db_name", "")

	username := ""
	pwd := ""
	readFromEnv := p.GetBool("db_read_creds_from_env", false)

	if readFromEnv {
		usernameEnvVariableKey := p.GetString("db_username_env", "")
		pwdEnvVariableKey := p.GetString("db_password_env", "")

		if usernameEnvVariableKey == "" {
			return fmt.Errorf("db username environment variable wasn't set")
		}

		if pwdEnvVariableKey == "" {
			return fmt.Errorf("db password environment variable wasn't set")
		}

		username = os.Getenv(usernameEnvVariableKey)
		pwd = os.Getenv(pwdEnvVariableKey)
	} else {
		username = p.GetString("db_username", "")
		pwd = p.GetString("db_password", "")
	}

	if username == "" {
		return fmt.Errorf("db username wasn't set")
	}
	if pwd == "" {
		return fmt.Errorf("db password wasn't set")
	}

	c.username = username
	c.pwd = pwd
	c.connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.username, c.pwd, c.host, c.port, c.dbName)
	return nil
}

func (c *EvaModuleRepositoryDatabaseConfig) GetConnectionString() string {
	return c.connectionString
}
