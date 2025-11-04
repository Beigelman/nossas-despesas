package config

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

func bindStructEnv(v *viper.Viper, obj any) error {
	// Get the reflection Value of the struct
	val := reflect.ValueOf(obj)
	// Check if the passed interface is a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("provided interface is not a struct")
	}

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		// Get the field type
		key := val.Type().Field(i).Tag.Get("mapstructure")
		if err := bindEnv(v, key); err != nil {
			return fmt.Errorf("config.bindStructEnv: %w", err)
		}
	}

	return nil
}

func bindEnv(v *viper.Viper, key string) error {
	if err := v.BindEnv(key); err != nil {
		return fmt.Errorf("config.bindEnv: %w", err)
	}

	return nil
}
