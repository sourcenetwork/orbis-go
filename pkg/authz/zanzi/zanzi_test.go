package zanzi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"

	"github.com/sourcenetwork/zanzi"
	"github.com/sourcenetwork/zanzi/pkg/api"
	"github.com/sourcenetwork/zanzi/pkg/domain"
	"github.com/sourcenetwork/zanzi/pkg/server"
)

func setupGRPC(address string) error {
	z, err := zanzi.New(
		zanzi.WithDefaultLogger(),
	)
	if err != nil {
		return err
	}

	server := server.NewServer(address)

	if err := server.Init(&z); err != nil {
		return err
	}

	go server.Run()
	return nil
}

func TestZanzi(t *testing.T) {
	port := rand.Int63n(55555) + 9999
	address := fmt.Sprintf("127.0.0.1:%d", port)
	err := setupGRPC(address)
	require.NoError(t, err)

	z, err := newGRPC(address)
	require.NoError(t, err)

	assert.Equal(t, "zanzi", z.Name())

	ctx := context.Background()
	require.NoError(t, setup(ctx, z))

	check, err := z.Check(ctx, "10/file:readme#read", "user:bob")
	require.NoError(t, err)
	require.True(t, check)
}

func setup(ctx context.Context, z *zanziGRPC) error {

	// Create Policy
	createResponse, err := z.policyClient.CreatePolicy(ctx, &api.CreatePolicyRequest{
		PolicyDefinition: &api.PolicyDefinition{
			PolicyYaml: policyYaml,
		},
	})
	if err != nil {
		return err
	}
	// log.Printf("Created policy: %v", createResponse.Record.Policy)

	// Set Relationships
	for _, relationship := range relationships {
		_, err = z.policyClient.SetRelationship(ctx, &api.SetRelationshipRequest{
			PolicyId:     createResponse.Record.Policy.Id,
			Relationship: &relationship,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

var policyYaml = `
id: 10
name: test
doc: test policy

resources:
  file:
    doc: file resource
    relations:
      owner:
        expr: _this
        types:
          - '*'
      parent:
        expr: _this
        types:
          - directory
      read:
        expr: owner + parent->read
        types: []
  directory:
    relations:
      owner:
        types:
          - user
          - group:member
          - group:owner
        expr: _this
      reader:
        expr: _this
        types:
          - user
          - group:member
      read:
        expr: owner + reader

  group:
    relations:
      member:
        expr: _this + owner
        types:
          - user
      owner:
        expr: _this
        types:
          - user
  user:

attributes:
  foo: bar
`

var builder domain.RelationshipBuilder

var relationships []domain.Relationship = []domain.Relationship{
	builder.Relationship("file", "readme", "owner", "user", "charlie"),
	builder.Relationship("file", "readme", "parent", "directory", "proj"),
	builder.EntitySet("directory", "proj", "owner", "group", "eng", "owner"),
	builder.EntitySet("directory", "proj", "reader", "group", "eng", "member"),
	builder.Relationship("group", "eng", "owner", "user", "alice"),
	builder.Relationship("group", "eng", "member", "user", "bob"),
}
