package cobracli

import (
	"fmt"
	"log"
	"path"
	"reflect"
	"strings"

	"github.com/sourcenetwork/orbis-go/config"

	"github.com/iancoleman/strcase"
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
	var err error

	var traverse func(v reflect.Value, path string)
	traverse = func(v reflect.Value, path string) {

		for i := 0; i < v.NumField(); i++ {

			field := v.Type().Field(i)
			name, tag := field.Name, field.Tag
			if path != "" {
				name = path + "." + name
			}

			f := v.Field(i)
			// TODO: is map handeled properly?
			if f.Kind().String() == "struct" {
				x := reflect.ValueOf(f.Interface())
				traverse(x, name)
				continue
			}

			snake := strcase.ToSnakeWithIgnore(name, ".")

			// Generate the Cobra command flag.
			// TODO: reflect on the type of the field and set the appropriate flag type.
			cmd.Flags().String(snake, tag.Get("default"), tag.Get("description"))

			// Bind the flag to the viper config.
			err = viper.BindPFlag(snake, cmd.Flags().Lookup(snake))
			if err != nil {
				return
			}
		}
	}

	x := reflect.ValueOf(cfg).Elem()
	traverse(x, "")

	return err
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

	if viper.ConfigFileUsed() != "" {
		log.Printf("using %s", viper.ConfigFileUsed())
	}

	// Unmarshal the config file into the config.Config struct.
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
