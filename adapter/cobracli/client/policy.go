package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

func PolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Interact with authorization policies",
	}

	cmd.AddCommand(
		DescribePolicyCmd(cfg),
		CreatePolicyCmd(cfg),
		RegisterPolicyCmd(cfg),
		SetRelationshipPolicyCmd(cfg),
		CheckPolicyCmd(cfg),
	)
	return cmd
}

func DescribePolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get and describe an existing policy",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hi from describe policy")
		},
	}
	return cmd
}

func CreatePolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new policy",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hi from create policy")
		},
	}
	return cmd
}

func RegisterPolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a resource instance in the policy",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hi from register policy")
		},
	}
	return cmd
}

func SetRelationshipPolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Create a relation",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hi from set policy")
		},
	}
	return cmd
}

func CheckPolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Evaluate a check call",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hi from check policy")
		},
	}
	return cmd
}
