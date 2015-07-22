package rpc

import (
	"golang.org/x/net/context"

	"github.com/robskillington/gokit-gorilla-mux-starter/deps"
	"github.com/robskillington/gokit-gorilla-mux-starter/models"
)

type CreateEntityRequest struct {
	Name string `json:"name"`
}

type CreateEntityResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type CreateEntity func(context.Context, *CreateEntityRequest) (*CreateEntityResponse, error)

func NewCreateEntity(injected *deps.All) CreateEntity {
	return func(ctx context.Context, req *CreateEntityRequest) (*CreateEntityResponse, error) {
		var entity *models.Entity
		var err error
		if entity, err = models.NewEntity(); err != nil {
			return nil, err
		}
		entity.Name = req.Name
		injected.EntityService.Save(entity)
		return &CreateEntityResponse{UUID: entity.UUID.String(), Name: entity.Name}, nil
	}
}
