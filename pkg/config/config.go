package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// If no arguments are passed by user, returns default config unless user defines it as environment variables
type Config struct {
	RedisAddr             string `arg:"env:REDIS_ADDRESS"`
	GlobalCacheExpiryTime int    `arg:"env:CACHE_EXPIRY_TIME"`
	CacheCapacity         int    `arg:"env:CACHE_CAPACITY"`
	ProxyPort             int    `arg:"env:PROXY_PORT"`
	MaxClients            int    `arg:"env:MAX_CLIENTS"`
}

func LookUpEnv(key string) (string, bool) {
	if r := goDotEnvVariable(key); r != "" {
		return r, true
	}
	return "", false

}
func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := LookUpEnv(key); ok {
		return val
	}
	return defaultVal
}
func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := LookUpEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func getConfig(fs *flag.FlagSet) []string {
	cfg := make([]string, 0, 10)
	fs.VisitAll(func(f *flag.Flag) {
		cfg = append(cfg, fmt.Sprintf("%s:%q", f.Name, f.Value.String()))
	})
	return cfg
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
func LoadConfig() *Config {
	config := &Config{}
	flag.StringVar(&config.RedisAddr, "redis-addr", LookupEnvOrString("REDIS_ADDRESS", "redis:6379"), "Address of backing Redis server")                                             // using default redis address port
	flag.IntVar(&config.GlobalCacheExpiryTime, "global-cache-expiry-time", LookupEnvOrInt("GLOBAL_CACHE_EXPIRY_TIME", 60*1000), "Cache expiry time(Please mention in milliseconds)") // todo for future: Handle time input as duration because user should not concern themselves with conversion
	flag.IntVar(&config.CacheCapacity, "cache-capacity", LookupEnvOrInt("CACHE_CAPACITY", 100), "Cache capacity")
	flag.IntVar(&config.ProxyPort, "proxy-port", LookupEnvOrInt("PROXY_PORT", 9000), "Port the proxy server listens on")
	flag.IntVar(&config.MaxClients, "max-clients", LookupEnvOrInt("MAX_CLIENTS", 1), "Concurrent client limit")
	flag.Parse()
	//log.Printf("app.config %v\n", getConfig(flag.CommandLine))
	return config

}
