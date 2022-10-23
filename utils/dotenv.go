package util

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func GodotEnv(key string) string {
	env := make(chan string, 1)

	if os.Getenv("GO_ENV") != "production" {
		_, b, _, _ := runtime.Caller(0)
		projectRootPath := filepath.Join(filepath.Dir(b), "../")
		err := godotenv.Load(projectRootPath + "/.env")
		if err != nil {
			panic(err)
		}
		env <- os.Getenv(key)
	} else {
		env <- os.Getenv(key)
	}

	return <-env
}
