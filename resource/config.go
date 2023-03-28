package resource

import (
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

type (
	DataConfig struct {
		Database *DatabaseConfig
		Redis    *RedisConfig
	}

	DatabaseConfig struct {
		Driver string
		Source string
	}

	RedisConfig struct {
		Network      string
		Addr         string
		ReadTimeout  *durationpb.Duration
		WriteTimeout *durationpb.Duration
	}
)
