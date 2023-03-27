package resource

import (
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

type (
	Config struct {
		Database *Database
		Redis    *Redis
	}

	Database struct {
		Driver string
		Source string
	}

	Redis struct {
		Network      string
		Addr         string
		ReadTimeout  *durationpb.Duration
		WriteTimeout *durationpb.Duration
	}
)
