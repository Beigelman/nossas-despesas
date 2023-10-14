package env

import (
	"fmt"
)

type ErrUndefinedEnvironment struct {
	env string
}

func (e ErrUndefinedEnvironment) Error() string {
	return fmt.Sprintf("undefined environment '%s'", e.env)
}

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
	return "", ErrUndefinedEnvironment{env}
}
