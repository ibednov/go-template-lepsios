package example

import (
	"github.com/ibednov/go-template-lepsios/internal/domain/example/service"
	"github.com/ibednov/go-template-lepsios/internal/domain/example/transport/api/http/v1/handler"
	"go.uber.org/fx"
)

var Module = fx.Module("domain.example",
	fx.Provide(
		service.New,
		handler.New,
	),
)
