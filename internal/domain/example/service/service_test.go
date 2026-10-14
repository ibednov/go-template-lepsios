package service_test

import (
	"context"
	"testing"

	"github.com/ibednov/go-template-lepsios/internal/domain/example/service"
)

func TestHello(t *testing.T) {
	svc := service.New()
	msg, err := svc.Hello(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if msg.Text == "" {
		t.Fatal("expected non-empty text")
	}
}
