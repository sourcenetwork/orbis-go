package app

import (
	"github.com/sourcenetwork/orbis-go/infra/logger"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

type Option func(a *App) error

func DefaultOptions() Option {
	return func(a *App) error {
		return nil
	}
}

func WithLogger(lg logger.Logger) Option {
	return func(a *App) error {
		a.lg = lg
		return nil
	}
}

func WithTransport(tp transport.Transport) Option {
	return func(a *App) error {
		a.tp = tp
		return nil
	}
}

// WithBulletin registers a BulletinBoard Service.
func WithBulletinService(bn bulletin.Bulletin) Option {
	return func(a *App) error {
		a.bn = bn
		return nil
	}
}

// WithSharing regisers a Proactive Sharing Service.
func WithSharingService(pss pss.PSS) Option {
	return func(a *App) error {
		a.pss = pss
		return nil
	}
}
