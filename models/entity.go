package models

import (
	uuid "github.com/nu7hatch/gouuid"
)

type Entity struct {
	UUID *uuid.UUID
}

func NewEntity() (*Entity, error) {
	var id *uuid.UUID
	var err error
	if id, err = uuid.NewV4(); err != nil {
		return nil, err
	}
	return &Entity{UUID: id}, nil
}
