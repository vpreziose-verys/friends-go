package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BethesdaNet/friends-go/cmd"
	"github.com/BethesdaNet/friends-go/internal/db/redis"
	"github.com/BethesdaNet/friends-go/internal/friends"
	"github.com/BethesdaNet/friends-go/internal/metric/relic"
	"github.com/BethesdaNet/friends-go/internal/provider"
)

var (
	name         = flag.String("name", "local-friends", "name of service")
	env          = flag.String("env", "local", "env of service")
	addr         = flag.String("addr", "localhost:10000", "address of the http listener")
	redisAddr    = flag.String("redisAddr", "localhost:6379", "address of redis")
	redisCluster = flag.Bool("redisCluster", true, "enable clustered redis agent")
	identityAddr = flag.String("identityAddr", "http://localhost:10001/identity", "address of identity service")
	identityKey  = flag.String("identityKey", "key-identity", "key for identity service")
	presenceAddr = flag.String("presenceAddr", "http://localhost:10001/presence", "address of presence service")
	presenceKey  = flag.String("presenceKey", "key-presence", "key for presence service")
	noteAddr     = flag.String("noteAddr", "http://localhost:10001/notification", "address of notification service")
	noteKey      = flag.String("noteKey", "key-note", "key for notification service")
	apmKey       = flag.String("apmKey", "", "apm agent license key")
	apmEnable    = flag.Bool("apmEnable", false, "enable apm agent")
	cpuprofile   = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile   = flag.String("memprofile", "", "write memory profile to file")
)

func main() {
	// trapfatal recovers thrown panics or fatal calls and evaulates them before
	// eventually panicing with more detailed stack data of the event.k
	defer trapfatal()

	// parse env vars and flags, refer to cmd pkg (cmd/env.go) on order of priority
	// and other functions to extract info about the running service
	cmd.ParseFlagsOrEnv()

	// create presence service configuration based cmd flags. Flag values are also
	// loaded by the above func (cmd.ParseFlagsOrEnv) if names match as env vars
	conf := friends.Config{
		Name: *name, Addr: *addr, Env: *env,
		Redis: redis.Config{Addr: strings.Split(*redisAddr, ","), Clustered: *redisCluster, TTL: -1, Retries: 3},
		Relic: relic.Config{Name: *name, Key: *apmKey, Enabled: *apmEnable},
	}

	// aws.Describe gathers important container, and aws-ecr meta data if available
	// only running where env conf.Env is not empty or local
	//aws.Describe(conf.Name, conf.Addr, conf.Env, conf.Meta)

	// populate provider map per subsystem (key) with addr and host. When operating
	// in the cloud these values are stored in AWS::SecretsManager by environment
	// and loaded via AWS::ECS TaskDef on spinup.
	conf.Provider = map[string]interface{}{
		"identity": &provider.IdentityConfig{Config: convert(*identityAddr, *identityKey), LookupURL: "/v2/lookup/identity/"},
		"presence": &provider.PresenceConfig{Config: convert(*presenceAddr, *presenceKey), PresenceURL: "/v1/presence", PresencePrivateURL: "/v1/presence-private"},
		"note":     &provider.NoteConfig{Config: convert(*noteAddr, *noteKey), AnnounceURL: "/v1/announcement", NoteURL: "/v1/notification"},
		"storage":  &provider.StorageConfig{},
	}

	// create redis agent and check for errors from the agent on init (ping) which
	// will attempt to retry (not redial) based on redis config
	dba, err := redis.Open(conf.Redis)
	check("dba: open:", err)

	// create relic agent and check for errors from the agent on init
	nra, err := relic.Open(conf.Relic)
	check("nra: open", err)

	// create presence service with config (require), redis agent (required), and
	// newrelic monitoring agent (not required) if not local env, service name, and
	// valid app auth keys populated via AWS::SecretsManager
	svc, err := friends.Open(conf, dba, nra)
	check("svc: open", err)

	r := svc.Routes()
	srv := &http.Server{Addr: *addr, Handler: r}

	// handle signals for handling teardown
	sigs, done := make(chan os.Signal, 2), make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Print("\r") // prevents ctrl+c terminal output

		// prepare ctx and defer fn (cancellation) at 1s on event of server stall/issue
		ctx, fn := context.WithTimeout(context.Background(), time.Second)
		defer fn()

		// call http server shutdown and force drain current connections with ctx
		srv.Shutdown(ctx)

		// after http service is shutdown continue to teardown other components of the
		// service cleanly.
		for _, stop := range []func(){svc.Close, dba.Close} {
			stop()
		}

		// close done channel / unblock to allow service close/exit.
		close(done)
	}()

	info("srv", fmt.Sprintf("ready > serving: %s%s", conf.Addr, conf.Target))

	// listen and serve the presence service
	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			info("srv: http", "closed")
		} else {
			info("srv: http error: %v", err)
		}
	}

	// done channel blocks only being closed when an valid os.Signal
	<-done

	info("srv", "closed")
}

// convert returns common provider config maintaining required addr and key attrs
func convert(addr, key string) provider.Config {
	return provider.Config{Addr: addr, Key: key}
}

// check evals error, if non-nil, it throws the fatal error up the stack
// to trapfatal, where it then calls log.Fatal
func check(what string, err error) {
	if err != nil {
		panic(fatal{what, err})
	}
}

// info is a convience func to log messages in a custom format / logger
func info(what, that interface{}) {
	log.Printf("%s: %v", what, that)
}

func path(uri string, private bool) string {
	out := fmt.Sprintf("/%s", uri)
	if private {
		out += "-private"
	}
	return out
}

type fatal struct {
	what string
	err  error
}

// trapfatal should be called once in a defer statement in the main function
func trapfatal() {
	err := recover()
	switch err := err.(type) {
	case nil:
		return
	case fatal:
		log.Fatalf("fatal: %s: %v", err.what, err.err)
	}
	panic(err)
}
