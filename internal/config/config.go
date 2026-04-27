package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr              string
	BalanceSystemAddr string
	DatabaseDsn       string
}

func NewConfig() Config {
	return Config{}
}

func (cfg *Config) InitConfig() Config {
	cfg.argsConfigurator()
	cfg.envConfigurator()

	return Config{
		Addr:              cfg.Addr,
		BalanceSystemAddr: cfg.BalanceSystemAddr,
		DatabaseDsn:       cfg.DatabaseDsn,
	}
}

func (cfg *Config) argsConfigurator() {
	addr := flag.String("a", "127.0.0.1:8080", "server ip:port")
	balanceSystemAddr := flag.String("r", "http://127.0.0.1:8080", "server result ip:port")
	databaseDsn := flag.String("d", "", "database connection string")

	flag.Parse()

	cfg.Addr = *addr
	cfg.BalanceSystemAddr = *balanceSystemAddr
	cfg.DatabaseDsn = *databaseDsn
}

func (cfg *Config) envConfigurator() {
	if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}
	if envDatabaseDsn := os.Getenv("DATABASE_URI"); envDatabaseDsn != "" {
		cfg.DatabaseDsn = envDatabaseDsn
	}
	if envBalanceSystemAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envBalanceSystemAddr != "" {
		cfg.BalanceSystemAddr = envBalanceSystemAddr
	}
}
