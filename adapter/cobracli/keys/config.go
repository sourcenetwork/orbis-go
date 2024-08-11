// Credit: This file is inspired by github.com/NathanBaulch/protoc-gen-cobra
package keys

import (
	"github.com/sourcenetwork/orbis-go/pkg/util/naming"

	"github.com/spf13/pflag"
)

type Config struct {
	KeyringBackend string // keyring backend
	KeyringPath    string // Keyring path argument for some backends
	KeyringService string // Keyring service argument for some backends

	UseEnvVars   bool
	EnvVarPrefix string

	EnvVarNamer naming.Namer
	FlagNamer   naming.Namer

	Output string
}

var DefaultConfig = &Config{
	KeyringBackend: "file",
	KeyringPath:    "$HOME/.orbis",
	UseEnvVars:     true,
	EnvVarPrefix:   "orbis",
	EnvVarNamer:    naming.UpperSnake,
	FlagNamer:      naming.LowerKebab,
	Output:         "yaml",
}

func (c *Config) BindFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&c.KeyringBackend, namer("KeyringBackend"), "k", c.KeyringBackend, "keyring backend to get identities from (file|os|test)")
	fs.StringVarP(&c.KeyringPath, namer("KeyringPath"), "", c.KeyringPath, "keyring path argument for some backends")
	fs.StringVarP(&c.KeyringService, namer("KeyringService"), "", c.KeyringService, "keyring service argument for some backends")
}

func namer(in string) string {
	return naming.LowerKebab(in)
}
