// Copyright © 2024 Kong Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultJWTSecret       = "kong"
	defaultJTWTokenTimeout = 30 * time.Minute
	defaultUsername        = "kong"
	defaultPassword        = "onward"
	defaultRequestTimeout  = 5 * time.Second
)

// Config is the configuration for the candidate take home exercise (SDET) to run.
type Config struct {
	// JWTSecret is the configuration for the secret key used for signing JWT tokens.
	JWTSecret string `yaml:"jwt_secret" mapstructure:"jwt_secret"`
	// JWTTokenTimeout is the timeout for token to expire.
	JWTTokenTimeout time.Duration `yaml:"jwt_token_timeout" mapstructure:"jwt_token_timeout"`
	// Username is the authorized user.
	Username string `yaml:"username" mapstructure:"username"`
	// Password is the password for the authorized user.
	Password string `yaml:"password" mapstructure:"password"`
	// RequestTimeout is the timeout for request operations.
	RequestTimeout time.Duration `yaml:"request_timeout" mapstructure:"request_timeout"`
}

// NewConfig creates a new configuration comprised of the configuration file,
// environment variables, and defaults.
func NewConfig() (*Config, error) {
	// Set default configuration vaules
	viper.SetDefault("jwt_secret", defaultJWTSecret)
	viper.SetDefault("username", defaultUsername)
	viper.SetDefault("password", defaultPassword)
	viper.SetDefault("request_timeout", defaultRequestTimeout)

	// Configuration setup for viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables to viper that do not have a corresponding default value
	viper.SetEnvPrefix("kong")

	// Enable automatic environment variable binding
	viper.AutomaticEnv()

	// Read in the configuration file and ignore not found errors as environment
	// variables will be used if the file is not found. If the required
	// configuration fields are not present then and error will be returned
	// further down the line.
	var config Config
	_ = viper.ReadInConfig()
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}
	return &config, nil
}
