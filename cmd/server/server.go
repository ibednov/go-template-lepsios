package server

import (
	"github.com/ibednov/go-template-lepsios/internal/app"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start HTTP API server",
		Run: func(cmd *cobra.Command, args []string) {
			app.Run()
		},
	}
}
