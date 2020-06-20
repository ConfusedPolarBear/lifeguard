// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"log"

	"github.com/spf13/viper"
)

func Load() {
	viper.SetConfigName("config")

	// TODO: add a proper path (like /etc/lifeguard) when packaging this for debian
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to load config: %s", err)
	}
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
	viper.WriteConfig()
}
