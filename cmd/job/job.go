package job

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "One-shot / background jobs",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "ping",
		Short: "No-op job stub",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("job ping ok")
			return nil
		},
	})
	return cmd
}
