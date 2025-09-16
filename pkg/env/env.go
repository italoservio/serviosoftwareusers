package env

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type Env struct {
	MONGODB_URI string
	PASS_SECRET string
	AUTH_SECRET string
}

func Load() *Env {
	env := &Env{}
	env.MONGODB_URI = os.Getenv("MONGODB_URI")
	env.PASS_SECRET = os.Getenv("PASS_SECRET")
	env.AUTH_SECRET = os.Getenv("AUTH_SECRET")

	if env.MONGODB_URI == "" {
		panic("MONGODB_URI is not set")
	}

	if env.PASS_SECRET == "" {
		panic("PASS_SECRET is not set")
	}

	if env.AUTH_SECRET == "" {
		panic("AUTH_SECRET is not set")
	}

	return env
}
