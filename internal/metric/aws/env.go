package aws

import (
	"os"
	"strings"
)

// envBlock is the process environment block's environment variables with
// lower case key values. Duplicate values are truncated in non-deterministic
// order because environment variables are case sensitive.
var envBlock = func() map[string]string {
	e := map[string]string{}
	for _, kv := range os.Environ() {
		k := strings.Split(kv, "=")[0]
		v := os.Getenv(k)
		e[strings.ToLower(k)] = v
	}
	return e
}()

func getenv(key string) string {
	return envBlock[strings.ToLower(key)]
}

func env() string {
	e := getenv("SERVICENAME")
	if e == "" {
		return "local"
	}
	s := strings.Split(e, "-")
	if len(s) == 0 {
		return e
	}
	return s[0]
}

func region() string {
	return getenv("REGION")
}

func hostname() string {
	return getenv("HOSTNAME")
}
