package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var Envs *Env

type Env struct {
	Port            int           `env:"PORT" envDefault:"8080"`
	AppEnv          string        `env:"APP_ENV" envDefault:"development"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"INFO"`
	AssetHost       string        `env:"ASSET_HOST" envDefault:"http://localhost/"`
	StorageProvider string        `env:"STORAGE_PROVIDER" envDefault:"local"`
	JWTSecret       string        `env:"JWT_SECRET" envDefault:"jwt secret"`
	TokenDuration   time.Duration `env:"TOKEN_DURATION" envDefault:"24h"`
	AllowOrigins    string        `env:"ALLOW_ORIGINS" envDefault:"*"`
	LocalFilePath   string        `env:"LocalFilePath" envDefault:"/tmp/sns-apigateway/"`
	MisesNodes      string        `env:"MISES_NODES" envDefault:"https://e1.mises.site:443,https://e2.mises.site:443,https://w1.mises.site:443,https://w2.mises.site:443"`
	RootPath        string
}

func init() {
	fmt.Println("apigateway env initializing...")
	_, b, _, _ := runtime.Caller(0)
	appEnv := os.Getenv("APP_ENV")
	projectRootPath := filepath.Dir(b) + "/../../"
	envPath := projectRootPath + ".env"
	appEnvPath := envPath + "." + appEnv
	localEnvPath := appEnvPath + ".local"
	_ = godotenv.Load(filtePath(localEnvPath, appEnvPath, envPath)...)
	Envs = &Env{}
	err := env.Parse(Envs)
	if err != nil {
		panic(err)
	}
	Envs.RootPath = projectRootPath
	fmt.Println("apigateway env loaded...")
}

func filtePath(paths ...string) []string {
	result := make([]string, 0)
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			result = append(result, path)
		}
	}
	return result
}
