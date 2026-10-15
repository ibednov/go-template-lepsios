package main

import (
	"github.com/ibednov/go-template-lepsios/cmd/job"
	"github.com/ibednov/go-template-lepsios/cmd/migrate"
	"github.com/ibednov/go-template-lepsios/cmd/server"
	"github.com/spf13/cobra"
)

func Execute() error {
	return NewRootCmd().Execute()
}

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "go-template-lepsios",
		Short: "go-template-lepsios CLI",
	}
	root.AddCommand(server.NewCmd())
	root.AddCommand(migrate.NewCmd())
	root.AddCommand(job.NewCmd())
	return root
}
