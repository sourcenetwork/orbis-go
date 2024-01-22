package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Config aggregates all the configuration options.
// It is data only, and minimal to none external dependencies.
// This implies only native types, and no external dependencies.
type Config struct {
	GRPC      GRPC
	Host      Host
	DKG       DKG
	Logger    Logger
	Transport Transport
	Bulletin  Bulletin
	DB        DB
	Authz     Authz
}

type Authz struct {
	Address string `default:"127.0.0.1:8080" description:"GRPC server address"`
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

type Transport struct {
}

type Bulletin struct {
	P2P struct {
	}
	SourceHub struct {
		AccountName   string `default:"alice" description:"Account name"`
		AddressPrefix string `default:"cosmos" description:"Address prefix"`
		Fees          string `default:"30stake" description:"Fees"`
		NodeAddress   string `default:"http://host.docker.internal:26657" description:"Node address"`
		RPCAddress    string `default:"tcp://host.docker.internal:26657" description:"RPC address"`
	}
}

type Host struct {
	Crypto struct {
		Type string `default:"ed25519" description:"crypto type"`
		Bits int    `default:"-1" description:"crypto bits, if selectable"`
		Seed int    `default:"0" description:"crypto seed"`
	}
	ListenAddresses []string `default:"/ip4/0.0.0.0/tcp/9000" description:"Host listen address string"`
	PersistentPeers []string `default:"" description:"Comma separated multiaddr strings of bootstrap peers."`
}

type DB struct {
	Path string `default:"data" description:"DB path"`
}

type configTypes interface {
	Host | DB | Bulletin | Transport | DKG | GRPC | Logger
}

func Default[T configTypes]() (T, error) {
	valT := new(T)

	x := reflect.ValueOf(valT).Elem()
	err := traverseAndBuildDefault(x)
	if err != nil {
		return *valT, fmt.Errorf("traverse: %w", err)
	}
	return *valT, nil
}

func traverseAndBuildDefault(v reflect.Value) error {
	// ensure struct
	for i := 0; i < v.NumField(); i++ {

		field := v.Type().Field(i)
		name, tag := field.Name, field.Tag

		f := v.Field(i)
		if !f.CanSet() {
			return fmt.Errorf("can't set field %s", name)
		}

		kind := f.Kind()

		// Generate the Cobra command flag.
		val, _ := tag.Get("default"), tag.Get("description")

		var err error
		var defaultValue any
		switch kind {
		case reflect.Struct:
			x := reflect.New(f.Type()).Elem()
			err := traverseAndBuildDefault(x)
			if err != nil {
				return fmt.Errorf("traverse: %w", err)
			}
			f.Set(x)
			continue
		case reflect.Bool:
			defaultValue, err = strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("parseBool: %q, %w", val, err)
			}
		case reflect.String:
			defaultValue = val
		case reflect.Int:
			defaultValue, err = strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("parseBool: %q, %w", val, err)
			}
		case reflect.Uint:
			defaultValue, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				return fmt.Errorf("parseBool: %q, %w", val, err)
			}
		case reflect.Float64:
			defaultValue, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return fmt.Errorf("parseBool: %q, %w", val, err)
			}
		case reflect.Slice:
			elmType := f.Type().Elem().Kind()
			if elmType != reflect.String {
				return fmt.Errorf("unsupported slice type: %q, for entry: %q", elmType, name)
			}
			defaultValue = strings.Split(val, ",")
		default:
			return fmt.Errorf("unsupported type: %q, for entry: %q", kind, name)
		}
		fmt.Println(name, kind, val, defaultValue)
		fmt.Println("can set", f.CanSet())
		f.Set(reflect.ValueOf(defaultValue))
	}

	return nil
}
