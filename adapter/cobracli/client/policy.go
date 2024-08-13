package client

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/sourcenetwork/orbis-go/adapter/cobracli"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/keyring"
	"github.com/sourcenetwork/zanzi/pkg/api"
	"github.com/sourcenetwork/zanzi/pkg/domain"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	flagCreateFile = "file"
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
		Use:   "describe policy-id",
		Short: "Get and describe an existing policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, ok := cobracli.FromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("couldn't get client context")
			}
			policyId := args[0]
			return RoundTrip(cmd.Context(), cfg, cfg.AuthzAddr, func(conn grpc.ClientConnInterface) error {
				policyClient := api.NewPolicyServiceClient(conn)

				getPolicyRequest := &api.GetPolicyRequest{
					Id: policyId,
				}
				resp, err := policyClient.GetPolicy(cmd.Context(), getPolicyRequest)
				if err != nil {
					return fmt.Errorf("get policy: %w", err)
				}

				ctx.Output().Print(resp.Record.Policy)
				return nil
			})
		},
	}
	return cmd
}

func CreatePolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [-f policy.yaml]",
		Short: "Create a new policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policyFilePath, err := cmd.Flags().GetString(flagCreateFile)
			if err != nil {
				return err
			}
			if policyFilePath == "" {
				return fmt.Errorf("policy.yaml path can't be empty")
			}
			policyYamlBuf, err := os.ReadFile(policyFilePath)
			if err != nil {
				return fmt.Errorf("can't read policy file: %w", err)
			}

			return RoundTrip(cmd.Context(), cfg, cfg.AuthzAddr, func(conn grpc.ClientConnInterface) error {
				policyClient := api.NewPolicyServiceClient(conn)

				policyCreateReq := &api.CreatePolicyRequest{
					PolicyDefinition: &api.PolicyDefinition{
						PolicyYaml: string(policyYamlBuf),
					},
				}
				resp, err := policyClient.CreatePolicy(cmd.Context(), policyCreateReq)
				if err != nil {
					return fmt.Errorf("create policy: %w", err)
				}

				fmt.Println("created policy:", resp.Record.Policy.Id)
				return nil
			})
		},
	}

	cmd.PersistentFlags().StringP(flagCreateFile, "f", "", "path to policy.yaml")
	cmd.MarkFlagRequired(flagCreateFile)
	return cmd
}

type setRelationOut struct {
	Overwritten bool
	PolicyId    string
	Resource    string
	Subject     string
	Relation    string
}

func RegisterPolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register policy-id resource-name resource-id",
		Short: "Register a resource instance in the policy",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, ok := cobracli.FromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("couldn't get client context")
			}
			did, err := fromDID(ctx.Keyring(), cfg.From)
			if err != nil {
				return fmt.Errorf("getting key DID identifier: %w", err)
			}

			// command args: <policyId> <resourceName> <resourceID>
			policyId := args[0]
			resourceName := args[1]
			resourceId := args[2]

			resp, err := doRelationshipRequest(cmd.Context(), cfg, policyId, resourceName,
				resourceId, did, "owner")

			if err != nil {
				return err
			}

			out := setRelationOut{
				Overwritten: resp.RecordOverwritten,
				PolicyId:    policyId,
				Resource:    resourceName + ":" + resourceId,
				Subject:     did,
				Relation:    "owner",
			}
			ctx.Output().Print(out)
			return nil
		},
	}
	return cmd
}

func SetRelationshipPolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set policy-id resource-name resource-id relation subject",
		Short: "Create a relation",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, ok := cobracli.FromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("unknown client context")
			}
			// command args: <policyId> <resourceName> <resourceID> <relation> <subject>
			policyId := args[0]
			resourceName := args[1]
			resourceId := args[2]
			relation := args[3]
			subject := args[4]

			resp, err := doRelationshipRequest(cmd.Context(), cfg, policyId, resourceName,
				resourceId, subject, relation)

			if err != nil {
				return err
			}

			out := setRelationOut{
				Overwritten: resp.RecordOverwritten,
				PolicyId:    policyId,
				Resource:    resourceName + ":" + resourceId,
				Subject:     subject,
				Relation:    relation,
			}
			ctx.Output().Print(out)
			return nil
		},
	}
	return cmd
}

func CheckPolicyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check policy-id subject permission",
		Short: "Evaluate a check call",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// command args: <policyId> <subject> <permission>
			policyId := args[0]
			subject := args[1]
			perm := args[2]

			parsedAccessReq, err := parsePermToAccessRequest(perm)
			if err != nil {
				return fmt.Errorf("parsing permission: %w", err)
			}
			parsedAccessReq.Subject = domain.NewEntity("user", subject)
			return RoundTrip(cmd.Context(), cfg, cfg.AuthzAddr, func(conn grpc.ClientConnInterface) error {
				relationClient := api.NewRelationGraphClient(conn)

				checkRequest := &api.CheckRequest{
					PolicyId:      policyId,
					AccessRequest: parsedAccessReq,
				}

				resp, err := relationClient.Check(cmd.Context(), checkRequest)
				if err != nil {
					return fmt.Errorf("failed check grpc: %w", err)
				}

				fmt.Println("valid:", resp.Result.Authorized)
				return nil
			})
		},
	}
	return cmd
}

var (
	permRegex = `^(?P<ResourceGroup>\w+):(?P<ResourceID>\w+)#(?P<Relation>\w+)$`
)

// permission is formatted as:
// PolicyID/ObjGroup:ObjID#relation
// we need to parse out:
// - ObjGroup
// - ObjID
// - relation
func parsePermToAccessRequest(permission string) (*domain.AccessRequest, error) {
	r, err := regexp.Compile(permRegex)
	if err != nil {
		return nil, err
	}

	if !r.Match([]byte(permission)) {
		return nil, fmt.Errorf("permission validation: %s", permission)
	}

	results := r.FindStringSubmatch(permission)
	if len(results) != 4 {
		return nil, fmt.Errorf("regex submatch size: %d", len(results))
	}
	return &domain.AccessRequest{
		Object:   domain.NewEntity(results[1], results[2]),
		Relation: results[3],
	}, nil
}

func createRelationship(policyId, resourceName, resourceId, subject, relation string) *domain.Relationship {
	return &domain.Relationship{
		Object: &domain.Entity{
			Resource: resourceName,
			Id:       resourceId,
		},
		Subject: &domain.Subject{
			Subject: &domain.Subject_Entity{
				Entity: &domain.Entity{
					Resource: "user",
					Id:       subject,
				},
			},
		},
		Relation: relation,
	}
}

func doRelationshipRequest(
	ctx context.Context,
	cfg *Config,
	policyId,
	resourceName,
	resourceId,
	subject,
	relation string) (*api.SetRelationshipResponse, error) {
	var resp *api.SetRelationshipResponse
	err := RoundTrip(ctx, cfg, cfg.AuthzAddr, func(conn grpc.ClientConnInterface) error {
		policyClient := api.NewPolicyServiceClient(conn)

		setRelationshipReq := &api.SetRelationshipRequest{
			PolicyId:     policyId,
			Relationship: createRelationship(policyId, resourceName, resourceId, subject, relation),
		}
		var err error
		resp, err = policyClient.SetRelationship(ctx, setRelationshipReq)
		return err
	})
	return resp, err
}

func fromDID(kr keyring.Keyring, from string) (string, error) {
	fromKey, err := kr.Get(from)
	if err != nil {
		return "", fmt.Errorf("getting key %s: %w", from, err)
	}
	if !crypto.IsPrivate(fromKey) {
		return "", fmt.Errorf("from key must be a private keypair")
	}
	fromKeyPriv := fromKey.(crypto.PrivateKey)
	did, err := fromKeyPriv.GetPublic().DID()
	if err != nil {
		return "", fmt.Errorf("getting key DID identifier: %w", err)
	}

	return did, nil
}
