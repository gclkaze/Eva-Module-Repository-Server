// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

func (c EvaModuleRepositoryDatabaseConfig) GetDBName() string {
	return c.dbName
}
