package service

import (
	"context"

	"github.com/ibednov/go-template-lepsios/internal/domain/example/entity"
	exampleerrors "github.com/ibednov/go-template-lepsios/internal/domain/example/errors"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Hello(_ context.Context) (entity.Message, error) {
	msg := entity.Message{Text: "hello from example domain"}
	if msg.Text == "" {
		return entity.Message{}, exampleerrors.ErrEmpty
	}
	return msg, nil
}
