package rpc

import (
	"golang.org/x/net/context"

	"github.com/robskillington/gokit-gorilla-mux-starter/deps"
	"github.com/robskillington/gokit-gorilla-mux-starter/models"
)

type CreateEntityRequest struct {
	UUID string `json:"uuid"`
}

type CreateEntityResponse struct {
	UUID string `json:"uuid"`
}

type CreateEntity func(context.Context, *CreateEntityRequest) (*CreateEntityResponse, error)

func NewCreateEntity(injected *deps.All) CreateEntity {
	return func(ctx context.Context, req *CreateEntityRequest) (*CreateEntityResponse, error) {
		var action *models.Entity
		var err error
		if action, err = models.NewEntity(); err != nil {
			return nil, err
		}
		injected.EntityService.Save(action)
		return &CreateEntityResponse{UUID: action.UUID.String()}, nil
	}
}
