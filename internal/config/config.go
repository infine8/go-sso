package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const CONFIG_PATH = "./config.yaml"

type Config struct {  
    Env            string     `yaml:"env" env-default:"local"`  
    StoragePath    string     `yaml:"storage_path" env-required:"true"`  
    GRPC           GRPCConfig `yaml:"grpc"`  
    MigrationsPath string  
    TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`  
}

type GRPCConfig struct {  
    Port    int           `yaml:"port"`  
    Timeout time.Duration `yaml:"timeout"`  
}

func MustLoad() *Config {  
    configPath := fetchConfigPath()  
    if configPath == "" {  
        panic("config path is empty") 
    }  

    return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {  
    // check if file exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        panic("config file does not exist: " + configPath)
    }

    var cfg Config

    if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
        panic("config path is empty: " + err.Error())
    }

    cfg.StoragePath = filepath.Join(rootPath(), cfg.StoragePath)

    return &cfg
}

func fetchConfigPath() string {
    return filepath.Join(rootPath(), CONFIG_PATH)
}

func rootPath() string {
    _, b, _, _ := runtime.Caller(0)
    return filepath.Join(filepath.Dir(b), "../..")
}