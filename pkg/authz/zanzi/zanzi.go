package zanzi

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sourcenetwork/orbis-go/pkg/authz"
	"github.com/sourcenetwork/zanzi/pkg/api"
	"github.com/sourcenetwork/zanzi/pkg/domain"
)

var (
	_ authz.Authz = (*zanziGRPC)(nil)
)

var (
	name      = "zanzi"
	permRegex = `^(?P<PolicyID>\w+)\/(?P<ResourceGroup>\w+):(?P<ResourceID>\w+)#(?P<Relation>\w+)$`
)

type zanziGRPC struct {
	conn           *grpc.ClientConn
	policyClient   api.PolicyServiceClient
	relationClient api.RelationGraphClient
}

func NewGRPC(address string) (*zanziGRPC, error) {
	cred := insecure.NewCredentials()
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(cred))
	if err != nil {
		return nil, err
	}

	return &zanziGRPC{
		conn:           conn,
		policyClient:   api.NewPolicyServiceClient(conn),
		relationClient: api.NewRelationGraphClient(conn),
	}, nil
}

func (z *zanziGRPC) Name() string {
	return name
}

func (z *zanziGRPC) Init(_ context.Context) error {
	return nil
}

func (z *zanziGRPC) Check(ctx context.Context, perm string, subject string) (bool, error) {
	checkReq, err := parsePermToCheckRequest(perm)
	if err != nil {
		return false, fmt.Errorf("parse permission: %w", err)
	}

	subjects := strings.SplitN(subject, ":", 2)
	if len(subjects) != 2 {
		return false, fmt.Errorf("subject validation: %s (%v size=%d)", subject, subjects, len(subjects))
	}
	checkReq.AccessRequest.Subject = domain.NewEntity(subjects[0], subjects[1])

	resp, err := z.relationClient.Check(ctx, checkReq)
	if err != nil {
		return false, fmt.Errorf("check rpc: %w", err)
	}

	return resp.Result.Authorized, nil
}

// permission is formatted as:
// PolicyID/ObjGroup:ObjID#relation
// we need to parse out:
// - PolicyID
// - ObjGroup
// - ObjID
// - relation
func parsePermToCheckRequest(permission string) (*api.CheckRequest, error) {
	r, err := regexp.Compile(permRegex)
	if err != nil {
		return nil, err
	}

	if !r.Match([]byte(permission)) {
		return nil, fmt.Errorf("permission validation: %s", permission)
	}

	results := r.FindStringSubmatch(permission)
	fmt.Println(results)
	if len(results) != 5 {
		return nil, fmt.Errorf("regex submatch size: %d", len(results))
	}
	return &api.CheckRequest{
		PolicyId: results[1],
		AccessRequest: &domain.AccessRequest{
			Object:   domain.NewEntity(results[2], results[3]),
			Relation: results[4],
		},
	}, nil
}
