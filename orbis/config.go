package orbis

type Config struct{}

type Option func(cfg *Config) error

func DefaultOptions() Option {
	panic("todo")
}

// WithBulletinService registers a BulletinBoard Service into the config
// func WithBulletinService(bulletin.WithProof) Option

// WithSharingService regisers a Proactive Sharing Service into the config
// func WithSharingService(pss.ProviderFn) Option
