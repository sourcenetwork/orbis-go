package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

func GetSecretClientCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("hi from secret get")
			return nil
		},
	}
	return cmd
}

func PutSecretClientCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put",
		Short: "Store a secret",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hi from secret put")
		},
	}
	return cmd
}
