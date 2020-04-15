package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// Env is the deployment environment the process is running in.
	Env = env()

	// Service is the identifier of the service set in ecs task definition. It does
	// this currently by parsing the SERVICENAME environment variable up to the
	// first dash.
	Service = service()
)

// Program returns the running exec name. This does not contain the directory.
func Program() string {
	p, _ := os.Executable()
	return filepath.Base(p)
}

// Args returns the parsed flags as they are interpreted after taking into
// account any environment variables injected into the process. The output
// represents the set of arguments that would be passed back into the process
// to achieve the same process behavior under those arguments.
func Args() []string {
	if !flag.Parsed() {
		ParseFlagsOrEnv()
	}
	set := []string{
		Program(),
	}
	flag.VisitAll(func(fl *flag.Flag) {
		if fl.Value == nil {
			return
		}
		set = append(set, "-"+fl.Name, fmt.Sprintf("%q", fl.Value))
	})
	set = append(set, "# GOVERSION", runtime.Version())
	set = append(set, "# SERVICE", service())
	set = append(set, "# ENV", env())
	return set
}

// ParseFlagsOrEnv parses the command line arguments. If the flag is not
// set, and an environment variable matching the flag name is set, ParseFlags
// assigns that environment variable's value to the flag.
//
// Order of priority (lower number is more authoritative):
// 	1) Flag
// 	2) Environment
// 	3) Default Value
func ParseFlagsOrEnv() {
	flag.VisitAll(func(fl *flag.Flag) {
		if env := getBlock(fl.Name); env != "" {
			fl.DefValue = env
			flag.Set(fl.Name, env)
		}
	})
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.Ltime)
}

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

func getBlock(key string) string {
	return envBlock[strings.ToLower(key)]
}

func env() string {
	return scan([]string{"env", "environment", "bnet_env"}, envBlock, "local")
}

func service() string {
	return scan([]string{"name", "servicename", "service_name"}, envBlock, "")
}

// scan takes in a set of keys and returns the first found by array index and if
// no env var is set the provided zero value will be used
func scan(keys []string, data map[string]string, zero string) string {
	for i, k := range keys {
		if v, ok := data[k]; ok && v != "" {
			log.Printf("scan: [%d/%d] %s=%s", i+1, len(keys), k, v)
			return v
		}
	}
	return zero
}
