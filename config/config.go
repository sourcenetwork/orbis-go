package config

// Config aggregates all the configuration options.
// It is data only, and minimal to none external dependencies.
// This implies only native types, and no external dependencies.
type Config struct {
	Logger Logger
	GRPC   GRPC
	Orbis  Orbis
	Ring   Ring
	Secret Secret
}

type Logger struct {
	Level  string `default:"DEBUG" description:"Log level"`
	Logger string `default:"zap" description:"Logger"`
	Zap    ZapLogger
}

type ZapLogger struct {
	Level    string `default:"debug" description:"Log level"`
	Encoding string `default:"dev" description:"Log encoding"`
}

type GRPC struct {
	GRPCURL string `default:"127.0.0.1:8080" description:"gRPC URL"`
	RESTURL string `default:"127.0.0.1:8090" description:"REST URL"`
}

type Orbis struct {
}

type Ring struct {
}

type Secret struct {
}
