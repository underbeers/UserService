package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	IsLocal   bool
	DB        *DB
	Listen    *Listen  `yaml:"listen"`
	Gateway   *Gateway `yaml:"gateway"`
	VersionDB int      `yaml:"db_version"`
}

type DB struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	NameDB     string `yaml:"name_db"`
	UserName   string `yaml:"user_name"`
	DBPassword string `env:"POSTGRES_PASSWORD" yaml:"password"`
}

type Listen struct {
	Port string `yaml:"port"`
	IP   string `yaml:"ip"`
}

type Gateway struct {
	Port  string `yaml:"port"`
	IP    string `yaml:"ip"`
	Label string `yaml:"label"`
}

type Service struct {
	Name      string `json:"name"`
	Port      string `json:"port"`
	IP        string `json:"ip"`
	Label     string `json:"label"`
	Endpoints []struct {
		URL       string   `json:"url"`
		Protected bool     `json:"protected"`
		Methods   []string `json:"methods"`
	} `json:"endpoints"`
}

func (db *DB) GetConnectionString() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		db.UserName, db.DBPassword, db.Host, db.Port, db.NameDB)
}

func GetConfig(isLocal bool) *Config {
	logger := log.Default()
	logger.Print("Read application configuration")
	instance := &Config{DB: &DB{}, IsLocal: isLocal}
	if err := cleanenv.ReadConfig("./conf/config.yml", instance); err != nil {
		help, _ := cleanenv.GetDescription(instance, nil)
		logger.Print(help)
		logger.Fatal(err)
	}

	confDBType := "docker"
	if instance.IsLocal {
		confDBType = "local"
	}
	if err := cleanenv.ReadConfig(fmt.Sprintf("./conf/db/%s.yml", confDBType), instance.DB); err != nil {
		help, _ := cleanenv.GetDescription(instance, nil)
		logger.Print(help)
		logger.Fatal(err)
	}

	err := cleanenv.ReadEnv(instance.DB)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return instance
}

func ReadConfig() *Config {
	instance := &Config{}
	err := cleanenv.ReadConfig("./conf/config.yml", instance)
	if err != nil {
		log.Fatalf("can't read config. %s", err.Error())
	}

	return instance
}

func ReadServicesList() *Service {
	instance := &Service{}
	err := cleanenv.ReadConfig("./conf/service.json", instance)
	if err != nil {
		log.Fatalf("can't read config. %s", err.Error())
	}

	return instance
}
