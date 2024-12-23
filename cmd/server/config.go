package main

import (
	"fmt"
	"log"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/peterbourgon/ff/v4/ffyaml"
)

type Environment string

const (
	Local Environment = "local"
)

func NewEnvironment(env string) Environment {
	switch env {
	case "local":
		return Local
	default:
		log.Printf("server: config: environment %q not available, using default %q\n", env, Local)
		return Local
	}
}

type Config struct {
	Environment Environment
	Server      Server
	Auth        Auth
	DB          DB
}

type Server struct {
	Host    string
	Port    string
	Timeout ServerTimeout
	Health  ServerHealth
}

type ServerTimeout struct {
	Read     int
	Write    int
	Idle     int
	Shutdown int
}

type ServerHealth struct {
	Timeout  int
	Cache    int
	Interval int
	Delay    int
	Retries  int
}

type Auth struct {
	JWT *JWT
}

type JWT struct {
	Algorithm string
	Key       string
}

type DB struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSL      string
	DSN      string
}

func NewConfig(args []string) (*Config, error) {
	fs := ff.NewFlagSet("okj")
	var (
		config                string
		env                   string
		serverHost            string
		serverPort            string
		serverTimeoutRead     int
		serverTimeoutWrite    int
		serverTimeoutIdle     int
		serverTimeoutShutdown int
		serverHealthTimeout   int
		serverHealthCache     int
		serverHealthInterval  int
		serverHealthDelay     int
		serverHealthRetries   int
		authJWTAlg            string
		authJWTKey            string
		dbHost                string
		dbPort                string
		dbName                string
		dbUser                string
		dbPasswd              string
		dbSSL                 string
	)
	fs.StringEnumVar(&config, 0, "config", "environment configuration file", "env.yaml")
	fs.StringEnumVar(&env, 0, "env", "build environment", "local")
	fs.StringVar(&serverHost, 0, "server.host", "", "server host address to listen for incoming requests")
	fs.StringVar(&serverPort, 0, "server.port", "", "server port number to listen for incoming requests")
	fs.IntVar(&serverTimeoutRead, 0, "server.timeout.read", 15, "number of seconds that the server will wait for reading requests")
	fs.IntVar(&serverTimeoutWrite, 0, "server.timeout.write", 15, "number of seconds that the server will wait for writing requests")
	fs.IntVar(&serverTimeoutIdle, 0, "server.timeout.idle", 60, "number of seconds that the server will wait for the next request")
	fs.IntVar(&serverTimeoutShutdown, 0, "server.timeout.shutdown", 30, "the duration for which the server gracefully wait for existing connections to finish")
	fs.IntVar(&serverHealthTimeout, 0, "server.health.timeout", 30, "if a single run of the check takes longer than timeout seconds then the check is considered to have failed")
	fs.IntVar(&serverHealthCache, 0, "server.health.cache", 5, "sets the duration for how long the aggregated health check result will be cached")
	fs.IntVar(&serverHealthInterval, 0, "server.health.interval", 30, "the health check will first run interval seconds after the program is started, and then again interval seconds after each previous check completes")
	fs.IntVar(&serverHealthDelay, 0, "server.health.delay", 5, "the initialization time for the program to bootstrap before the health check begins")
	fs.IntVar(&serverHealthRetries, 0, "server.health.retries", 3, "the number of consecutive failures of the health check for the container to be considered unhealthy")
	fs.StringVar(&authJWTAlg, 0, "auth.jwt.alg", "", "algorithm that was used for signing the JWT token")
	fs.StringVar(&authJWTKey, 0, "auth.jwt.key", "", "key that was used for signing the JWT token")
	fs.StringVar(&dbHost, 0, "db.host", "", "database host address")
	fs.StringVar(&dbPort, 0, "db.port", "", "database port number")
	fs.StringVar(&dbName, 0, "db.name", "", "database name")
	fs.StringVar(&dbUser, 0, "db.user", "", "database user")
	fs.StringVar(&dbPasswd, 0, "db.passwd", "", "database password")
	fs.StringVar(&dbSSL, 0, "db.ssl", "", "database ssl mode")

	if err := ff.Parse(fs, args[1:],
		ff.WithEnvVarPrefix("OKJ"),
		ff.WithConfigFileParser(ffyaml.Parse),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigIgnoreUndefinedFlags(),
	); err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		fmt.Printf("ERROR\n%v\n", err)
		return &Config{}, err
	}

	return &Config{
		Environment: NewEnvironment(env),
		Server: Server{
			Host: serverHost,
			Port: serverPort,
			Timeout: ServerTimeout{
				Read:     serverTimeoutRead,
				Write:    serverTimeoutWrite,
				Idle:     serverTimeoutIdle,
				Shutdown: serverTimeoutShutdown,
			},
			Health: ServerHealth{
				Timeout:  serverHealthTimeout,
				Cache:    serverHealthCache,
				Interval: serverHealthInterval,
				Delay:    serverHealthDelay,
				Retries:  serverHealthRetries,
			},
		},
		Auth: Auth{
			&JWT{
				Algorithm: authJWTAlg,
				Key:       authJWTKey,
			},
		},
		DB: DB{
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
			User:     dbUser,
			Password: dbPasswd,
			SSL:      dbSSL,
		},
	}, nil
}
