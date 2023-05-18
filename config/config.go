package config

// Config aggregates all the configuration options.
// It is data only, and minimal to none external dependencies.
// This implies only native types, and no external dependencies.
type Config struct {
	GRPC      GRPC
	Host      Host
	DKG       DKG
	Logger    Logger
	Ring      Ring
	Secret    Secret
	Transport Transport
	Bulletin  Bulletin
	DB        DB
}

type Logger struct {
	Level  string `default:"debug" description:"Log level"`
	Logger string `default:"zap" description:"Logger"`
	Zap    struct {
		Encoding string `default:"dev" description:"Log encoding"`
	}
}

type GRPC struct {
	GRPCURL string `default:"127.0.0.1:8080" description:"gRPC URL"`
	RESTURL string `default:"127.0.0.1:8090" description:"REST URL"`
	Logging bool   `default:"false" description:"debug mode"`
}

type DKG struct {
	Repo      string `default:"simpledb" description:"DKG repo"`
	Transport string `default:"p2ptp" description:"DKG transport"`
	Bulletin  string `default:"p2pbb" description:"DKG Bulletin"`
}

type Ring struct {
}

type Secret struct {
}

type Transport struct {
	Rendezvous string `default:"orbis-transport" description:"Rendezvous string"`
}

type Bulletin struct {
	Rendezvous string `default:"orbis-bulletin" description:"Rendezvous string"`
}

type Host struct {
	Crypto struct {
		Type string `default:"ed25519" description:"crypto type"`
		Bits int    `default:"-1" description:"crypto bits, if selectable"`
		Seed int    `default:"0" description:"crypto seed"`
	}
	ListenAddresses []string `default:"/ip4/0.0.0.0/tcp/9000" description:"Host listen address string"`
	BootstrapPeers  []string `mapstructure:"bootstrap_peers" default:"" description:"Comma separated multiaddr strings of bootstrap peers. If empty, the node will run in bootstrap mode"`
	// Rendezvous      string   `default:"orbis" description:"Rendezvous string"`
}

type DB struct {
	Path string `default:"data" description:"DB path"`
}
