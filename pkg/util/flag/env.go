// Credit: This file is copied from github.com/NathanBaulch/protoc-gen-cobra/naming
// It therefore retains the original copyright and license.
package flag

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sourcenetwork/orbis-go/pkg/util/naming"

	"github.com/spf13/pflag"
)

func SetFlagsFromEnv(fs *pflag.FlagSet, search bool, namer naming.Namer, prefixes ...string) (err error) {
	parts := make([]string, 0, len(prefixes)+1)
	for _, prefix := range prefixes {
		if prefix != "" {
			parts = append(parts, namer(prefix))
		}
	}

	fs.VisitAll(func(f *pflag.Flag) {
		if err != nil || f.Changed {
			return
		}
		if search && len(parts) > 0 {
			for i := len(parts); i > 0; i-- {
				if err = setFlagFromEnv(namer, parts[:i], f); err == errNotFound {
					err = nil
				} else {
					return
				}
			}
		} else if err = setFlagFromEnv(namer, parts, f); err == errNotFound {
			err = nil
		}
	})

	return
}

var errNotFound = errors.New("not found")

func setFlagFromEnv(namer naming.Namer, parts []string, f *pflag.Flag) error {
	name := strings.Join(append(parts, namer(f.Name)), "_")
	if val := os.Getenv(name); val != "" {
		if err := f.Value.Set(val); err != nil {
			return fmt.Errorf("environment variable %s: %v", name, err)
		}
		return nil
	}
	return errNotFound
}
