package cobracli

import (
	"fmt"
	"log"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/sourcenetwork/orbis-go/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "orbis" // Without extension.
	envPrefix         = "ORBIS" // Without leading underscore.
)

// StartCmd returns a Cobra command for starting the Orbis server.
func StartCmd(setup func(config.Config) error) (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start Orbis server",
		RunE: func(cmd *cobra.Command, args []string) error {

			file := cmd.Flag("config").Value.String()
			cfg, err := readConfigFile(file)
			if err != nil {
				return fmt.Errorf("read config file: %w", err)
			}

			return setup(cfg)
		},
	}

	cmd.Flags().String("config", defaultConfigName+".yaml", "Config filename")

	// builds command flags from the config.Config struct.
	err := buildCmdFlags(cmd)
	if err != nil {
		return nil, fmt.Errorf("build command flags: %w", err)
	}

	return cmd, nil
}

// buildCmdFlags builds command flags from the config.Config struct.
// TODO: the reflection code is a bit messy. Refactor.
func buildCmdFlags(cmd *cobra.Command) error {

	cfg := &config.Config{}

	var traverse func(v reflect.Value, path string) error
	traverse = func(v reflect.Value, path string) error {

		for i := 0; i < v.NumField(); i++ {

			field := v.Type().Field(i)
			name, tag := field.Name, field.Tag
			if path != "" {
				name = path + "." + name
			}

			f := v.Field(i)
			kind := f.Kind()

			snake := toSnakeCase(name)

			// Generate the Cobra command flag.
			val, desc := tag.Get("default"), tag.Get("description")

			switch kind {
			case reflect.Struct:
				x := reflect.ValueOf(f.Interface())
				err := traverse(x, name)
				if err != nil {
					return fmt.Errorf("traverse: %w", err)
				}
				continue
			case reflect.Bool:
				parsed, err := strconv.ParseBool(val)
				if err != nil {
					return fmt.Errorf("parseBool: %q, %w", val, err)
				}
				cmd.Flags().Bool(snake, parsed, desc)
			case reflect.String:
				cmd.Flags().String(snake, val, desc)
			case reflect.Int:
				parsed, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("parseBool: %q, %w", val, err)
				}
				cmd.Flags().Int(snake, parsed, desc)
			case reflect.Uint:
				parsed, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					return fmt.Errorf("parseBool: %q, %w", val, err)
				}
				cmd.Flags().Uint(snake, uint(parsed), desc)
			case reflect.Float64:
				parsed, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return fmt.Errorf("parseBool: %q, %w", val, err)
				}
				cmd.Flags().Float64(snake, parsed, desc)
			case reflect.Slice:
				// TODO: support other slice types.
				elmType := f.Type().Elem().Kind()
				if elmType != reflect.String {
					return fmt.Errorf("unsupported slice type: %q, for entry: %q", elmType, name)
				}
				cmd.Flags().StringSlice(snake, strings.Split(val, ","), desc)
			default:
				return fmt.Errorf("unsupported type: %q, for entry: %q", kind, name)
			}

			// Bind the flag to the viper config.
			err := viper.BindPFlag(snake, cmd.Flags().Lookup(snake))
			if err != nil {
				return fmt.Errorf("bind flag: %w", err)
			}
		}
		return nil
	}

	x := reflect.ValueOf(cfg).Elem()
	err := traverse(x, "")
	if err != nil {
		return fmt.Errorf("traverse: %w", err)
	}

	return nil
}

func readConfigFile(file string) (config.Config, error) {

	// Environment variables can't have dashes and dots.
	// Replace them with underscores. e.g. --grpc.grpcurl to GRPC_GRPCURL
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set the prefix for environment variables.
	viper.SetEnvPrefix(envPrefix)

	// Read in environment variables that match. e.g. ORBIS_GRPC_GRPCURL.
	viper.AutomaticEnv()

	// Set the config name (without extension).
	viper.SetConfigName(defaultConfigName)

	// Search config in current directory with name "orbis" (without extension).
	viper.AddConfigPath(".")

	// Search config in home directory with name ".orbis" (without extension).
	viper.AddConfigPath(path.Join("$HOME", "."+defaultConfigName))

	if file != defaultConfigName+".yaml" {
		// Use the given config file without searching.
		viper.SetConfigFile(file)
	}

	var cfg config.Config

	err := viper.ReadInConfig()
	if err != nil {
		// Config file is optional.
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return cfg, fmt.Errorf("read config file: %w", err)
		}
		log.Printf("no config file found")
	}

	// Unmarshal the config file into the config.Config struct.
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func toSnakeCase(str string) string {
	var matchAllCap = regexp.MustCompile("([a-z])([A-Z])")
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}
