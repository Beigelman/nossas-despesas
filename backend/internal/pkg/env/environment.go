package env

import (
	"fmt"
)

type Environment string

func (e Environment) String() string {
	return string(e)
}

const (
	Production  Environment = "prd"
	Development Environment = "development"
	Staging     Environment = "stg"
)

func Parse(env string) (Environment, error) {
	switch env {
	case "prd", "production":
		return Production, nil
	case "stg", "staging":
		return Staging, nil
	case "dev", "development":
		return Development, nil
	}
	return "", fmt.Errorf("undefined environment: %s", env)
}

func MustParse(env string) Environment {
	parsed, err := Parse(env)
	if err != nil {
		panic(err)
	}

	return parsed
}
