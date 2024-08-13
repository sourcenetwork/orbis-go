// Credit: This package is inspired by github.com/NathanBaulch/protoc-gen-cobra
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sourcenetwork/orbis-go/adapter/cobracli/keys"
	"github.com/sourcenetwork/orbis-go/pkg/util/naming"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Keyring *keys.Config

	From   string // keyring identity to execute with
	RingId string // secret ring ID

	ServerAddr   string        // remote orbis server address
	AuthzAddr    string        // remote authz server address
	Timeout      time.Duration // grpc request timeout
	UseEnvVars   bool
	EnvVarPrefix string

	EnvVarNamer naming.Namer
	FlagNamer   naming.Namer

	// TLS config
	TLS                bool
	ServerName         string
	InsecureSkipVerify bool
	CACertFile         string
	CertFile           string
	KeyFile            string

	Output string
}

var DefaultConfig = &Config{
	Keyring: keys.DefaultConfig,

	ServerAddr:   "localhost:8081",
	AuthzAddr:    "localhost:8080",
	Timeout:      10 * time.Second,
	UseEnvVars:   true,
	EnvVarPrefix: "orbis",
	EnvVarNamer:  naming.UpperSnake,
	FlagNamer:    naming.LowerKebab,
	Output:       "yaml",
}

func (c *Config) BindFlags(fs *pflag.FlagSet) {
	c.Keyring.BindFlags(fs)
	fs.StringVarP(&c.From, namer("From"), "", c.From, "keyring identity to use")
	fs.StringVarP(&c.RingId, namer("RingId"), "", c.RingId, "Secret Ring ID")
	fs.StringVarP(&c.ServerAddr, namer("ServerAddr"), "s", c.ServerAddr, "orbis service address in the form host:port")
	fs.StringVarP(&c.ServerAddr, namer("AuthzAddr"), "z", c.AuthzAddr, "authorization service address in the form host:port")
	fs.DurationVar(&c.Timeout, namer("Timeout"), c.Timeout, "client connection timeout")
	fs.BoolVar(&c.TLS, namer("TLS"), c.TLS, "enable TLS")
	fs.StringVar(&c.ServerName, namer("TLS ServerName"), c.ServerName, "TLS server name override")
	fs.BoolVar(&c.InsecureSkipVerify, namer("TLS InsecureSkipVerify"), c.InsecureSkipVerify, "INSECURE: skip TLS checks")
	fs.StringVar(&c.CACertFile, namer("TLS CACertFile"), c.CACertFile, "CA certificate file")
	fs.StringVar(&c.CertFile, namer("TLS CertFile"), c.CertFile, "client certificate file")
	fs.StringVar(&c.KeyFile, namer("TLS KeyFile"), c.KeyFile, "client key file")
	fs.StringVarP(&c.Output, namer("Output"), "", c.Output, "output format (text|json|yaml)")
}

func (c *Config) dialOpts(ctx context.Context, opts *[]grpc.DialOption) error {
	if c.TLS {
		tlsConfig := &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify}
		if c.CACertFile != "" {
			caCert, err := os.ReadFile(c.CACertFile)
			if err != nil {
				return fmt.Errorf("ca cert: %v", err)
			}
			certPool := x509.NewCertPool()
			certPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = certPool
		}
		if c.CertFile != "" {
			if c.KeyFile == "" {
				return fmt.Errorf("key file not specified")
			}
			pair, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
			if err != nil {
				return fmt.Errorf("cert/key: %v", err)
			}
			tlsConfig.Certificates = []tls.Certificate{pair}
		}
		if c.ServerName != "" {
			tlsConfig.ServerName = c.ServerName
		} else {
			addr, _, _ := net.SplitHostPort(c.ServerAddr)
			tlsConfig.ServerName = addr
		}
		cred := credentials.NewTLS(tlsConfig)
		*opts = append(*opts, grpc.WithTransportCredentials(cred))
	} else {
		*opts = append(*opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// TODO MIGHT BE NECESSARY FOR JWS Header Auth
	//
	// for _, dialer := range c.preDialers {
	// 	if err := dialer(ctx, opts); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func RoundTrip(ctx context.Context, cfg *Config, addr string, fn func(conn grpc.ClientConnInterface) error) error {
	var err error

	opts := []grpc.DialOption{grpc.WithBlock()}
	if err := cfg.dialOpts(ctx, &opts); err != nil {
		return err
	}

	if cfg.Timeout > 0 {
		var done context.CancelFunc
		ctx, done = context.WithTimeout(ctx, cfg.Timeout)
		defer done()
	}

	cc, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("timeout dialing server: %s", addr)
		}
		return err
	}
	defer cc.Close()

	return fn(cc)
}

func namer(in string) string {
	return naming.LowerKebab(in)
}
