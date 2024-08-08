package grpcserver

import (
	"context"
	"fmt"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/app"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/transport"

	icpb "github.com/libp2p/go-libp2p/core/crypto/pb"
)

type transportService struct {
	transportv1alpha1.UnimplementedTransportServiceServer
	app *app.App
}

func newTransportService(app *app.App) *transportService {
	return &transportService{
		app: app,
	}
}

func (s *transportService) GetHost(ctx context.Context, req *transportv1alpha1.GetHostRequest) (*transportv1alpha1.GetHostResponse, error) {

	tp, err := do.InvokeNamed[transport.Transport](s.app.Injector(), req.Transport)
	if err != nil {
		return nil, fmt.Errorf("not found")
	}

	raw, err := tp.Host().PublicKey().Raw()
	if err != nil {
		return nil, err
	}

	kt := icpb.KeyType(tp.Host().PublicKey().Type())
	resp := &transportv1alpha1.GetHostResponse{
		Node: &transportv1alpha1.Node{
			Id:      tp.Host().ID(),
			Address: tp.Host().Address().String(),
			PublicKey: &icpb.PublicKey{
				Type: &kt,
				Data: raw,
			},
		},
	}

	return resp, nil
}

func (s *transportService) GetPeers(ctx context.Context, req *transportv1alpha1.GetPeersRequest) (*transportv1alpha1.GetPeersResponse, error) {
	var peers []*transportv1alpha1.Node
	peerstore := s.app.Host().Network().Peerstore()
	for _, p := range s.app.Host().Network().Peers() {
		addr := peerstore.PeerInfo(p)
		peers = append(peers, &transportv1alpha1.Node{
			Id:      addr.ID.String(),
			Address: addr.String(),
		})
	}

	return &transportv1alpha1.GetPeersResponse{
		Nodes: peers,
	}, nil
}
