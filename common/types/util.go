package types

import (
	"fmt"
	"os"
	"strings"
)

// GetEnvBool returns bool value from env if var is set to valid bool.
// If var is not set then false will be returned. An error will only be
// returned of the env var is not a type bool.
func GetEnvBool(env string) (bool, error) {

	s := os.Getenv(env)
	if s == "" {
		return false, nil
	}

	switch strings.ToUpper(s) {

	case "TRUE":
		return true, nil

	case "FALSE":
		return false, nil

	}

	return false, fmt.Errorf("env variable %s is invalid. It should be either true or false", env)
}
