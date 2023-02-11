package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	DebugMode bool
	DB        *DB
	Listen    *Listen  `yaml:"listen"`
	Gateway   *Gateway `yaml:"gateway"`
	VersionDB int      `yaml:"db_version"`
	EmailFrom string   `yaml:"EMAIL_FROM"`
	SMTPHost  string   `yaml:"SMTP_HOST"`
	SMTPPass  string   `yaml:"SMTP_PASS"`
	SMTPPort  int      `yaml:"SMTP_PORT"`
	SMTPUser  string   `yaml:"SMTP_USER"`
}

type DB struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	NameDB     string `yaml:"name_db"`
	UserName   string `yaml:"user_name"`
	DBPassword string `yaml:"password"`
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

func GetConfig(debugMode bool) *Config {
	logger := log.Default()
	//logger.Print("Read application configuration")
	instance := &Config{DB: &DB{}, DebugMode: debugMode}
	if err := cleanenv.ReadConfig("./conf/config.yml", instance); err != nil {
		help, _ := cleanenv.GetDescription(instance, nil)
		logger.Print(help)
		logger.Fatal(err)
	}

	if debugMode {
		dbConfigName := "DBConfig"
		if err := cleanenv.ReadConfig(fmt.Sprintf("./conf/db/%s.yml", dbConfigName), instance.DB); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Print(help)
			logger.Fatal(err)
		}
	} else {
		instance.Gateway = &Gateway{
			IP:    getEnv("GATEWAY_IP", ""),
			Port:  getEnv("GATEWAY_PORT", ""),
			Label: getEnv("GATEWAY_LABEL", ""),
		}

		instance.Listen.IP = "127.0.0.1"
		instance.Listen.Port = "6001"

		instance.DB = &DB{
			Host:       getEnv("POSTGRES_HOST", ""),
			Port:       getEnv("POSTGRES_PORT", ""),
			NameDB:     getEnv("POSTGRES_DB_NAME", ""),
			UserName:   getEnv("POSTGRES_USER", ""),
			DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		}
	}

	return instance
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
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
