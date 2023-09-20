package util

import (
	"github.com/spf13/cobra"
)

func UtilCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "util",
		Short: "util",
		Long:  `a set of utility commands for working with cryptographic primitives.`,
	}

	cmd.AddCommand(encryptCmd())
	cmd.AddCommand(decryptCmd())
	cmd.AddCommand(createKeypair())
	cmd.AddCommand(createDIDCmd())
	cmd.AddCommand(createJWTCmd())

	return cmd
}
