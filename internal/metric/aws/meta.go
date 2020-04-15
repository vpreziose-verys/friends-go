package aws

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/BethesdaNet/friends-go/internal/db/ecr"
)

// AWS: EC2 Instance Metadata
// The link below explains in detail how to retrieve specific meta data
// ---------
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html

type MetaType string

const (
	MetaPattern    string   = "http://169.254.169.254/latest/meta-data/%v"
	MetaStats      MetaType = "stats"
	MetaIPV4                = "local-ipv4"
	MetaInstanceHN          = "hostname"
	MetaInstanceID          = "instance-id"
	MetaInstanceIP          = "instance-ip"
	MetaInstanceAZ          = "placement/availability-zone/"
)

var (
	defaultDead   = time.Duration(time.Second * 3)
	defaultSnooze = time.Duration(defaultDead / 2)
	defaultStats  = []MetaType{
		MetaInstanceHN,
		MetaInstanceID,
		MetaInstanceIP,
		MetaInstanceAZ,
		MetaStats,
	}
)

// Describe populates the Meta struct with aws / docker specific information
func Describe(name, addr, env string, m *Meta) {
	_, port, _ := net.SplitHostPort(addr)
	m = &Meta{
		Name:   name,
		Port:   port,
		Env:    env,
		Host:   hostname(),
		Region: region(),
		PID:    os.Getpid(),
	}
	if m.Env != "" && m.Env != "local" && len(defaultStats) > 0 {
		for _, f := range defaultStats {
			ctx, cancel := context.WithTimeout(context.Background(), defaultDead)
			defer cancel()
			go func() {
				select {
				case <-time.After(defaultSnooze):
					log.Printf("meta: request(%s): context sleepover: %d (next: deadline=%d)", f, defaultSnooze, defaultDead)
				case <-ctx.Done():
					log.Printf("meta: request(%s): context deadline:  %d (err: %v)", f, defaultDead, ctx.Err())
					return
				}
			}()
			m.getMeta(ctx, f)
		}
	}
}

// Meta maintains attributes from the container used for platform and logging
type Meta struct {
	Name       string  `json:"api_name"`
	Env        string  `json:"env"`
	Host       string  `json:"hostname"`
	PID        int     `json:"process_id"`
	Port       string  `json:"port"`
	Region     string  `json:"region"`
	InstanceNM *string `json:"instance_name,omitempty"`
	InstanceID *string `json:"instance_id,omitempty"`
	InstanceIP *string `json:"instance_ip,omitempty"`
	InstanceAZ *string `json:"instance_az,omitempty"`

	container *ecr.Container
}

func (m *Meta) getMeta(ctx context.Context, mt MetaType) {
	switch mt {
	case MetaInstanceHN:
		if v, err := inspect(ctx, MetaInstanceHN); err == nil {
			m.InstanceNM = &v
		}
	case MetaInstanceID:
		if v, err := inspect(ctx, MetaInstanceID); err == nil {
			m.InstanceID = &v
		}
	case MetaInstanceAZ:
		if v, err := inspect(ctx, MetaInstanceAZ); err == nil {
			m.InstanceAZ = &v
		}
	case MetaInstanceIP:
		if v, err := inspect(ctx, MetaIPV4); err == nil {
			m.InstanceIP = &v // MetaIPV4 & MetaInstanceIP?
		}
	case MetaStats:
		if c, err := ecr.Stats(); err == nil {
			m.container = c
		}
	}
}

func inspect(ctx context.Context, mt MetaType) (string, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf(MetaPattern, mt), nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return string(data), err
}
