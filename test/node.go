package test

import (
	"testing"

	"github.com/OpenBazaar/openbazaar-go/core"
	"github.com/OpenBazaar/openbazaar-go/ipfs"
)

func NewNode(t *testing.T) *core.OpenBazaarNode {
	repo, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	ipfsNode, err := ipfs.NewMockNode()
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := ipfs.MockCmdsCtx()
	if err != nil {
		t.Fatal(err)
	}

	return &core.OpenBazaarNode{
		Context:   ctx,
		RepoPath:  GetRepoPath(),
		IpfsNode:  ipfsNode,
		Datastore: repo.DB,
	}
}
