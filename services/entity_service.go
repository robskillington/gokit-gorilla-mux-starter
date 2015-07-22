package services

import (
	"fmt"

	"github.com/robskillington/gokit-gorilla-mux-starter/models"
)

type EntityService struct{}

func (s *EntityService) Save(action *models.Entity) error {
	fmt.Printf("Saving %s\n", action.UUID.String())
	return nil
}
