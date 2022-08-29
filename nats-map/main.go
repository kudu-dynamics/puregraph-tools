/*
 * nats-map
 * ========
 *
 * Utility to dump messages from NATS/STAN to stdout.
 */
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"github.com/spf13/viper"
)

type config struct {
	LastProcessed   uint64
	Mapper          string `mapstructure:"mapper"`
	NatsUrl         string `mapstructure:"nats_url"`
	StanChannel     string `mapstructure:"stan_channel"`
	StanClientId    string `mapstructure:"stan_client_id"`
	StanCluster     string `mapstructure:"stan_cluster"`
	StanDurableName string `mapstructure:"stan_durable_name" optional:"true"`
	StanQueueGroup  string `mapstructure:"stan_queue_group" optional:"true"`
}

const usageStr = `
Usage: nats-map

Env:
    MAPPER             Path to script to map over messages (input on stdin)
    NATS_URL           NATS Streaming server URL
    STAN_CHANNEL       NATS Streaming channel
    STAN_CLIENT_ID     NATS Streaming client ID
    STAN_CLUSTER       NATS Streaming cluster name
    STAN_DURABLE_NAME  NATS Streaming durable subscriber name
    STAN_QUEUE_GROUP   NATS Streaming queue group
`

func loadcfg() (cfg config, err error) {
	cfg.LastProcessed = 0

	v := viper.New()
	v.SetDefault("mapper", "")
	v.SetDefault("nats_url", "nats://localhost:4222")
	v.SetDefault("stan_channel", "")
	v.SetDefault("stan_client_id", "")
	v.SetDefault("stan_cluster", "")
	v.SetDefault("stan_durable_name", "")
	v.SetDefault("stan_queue_group", "")
	v.AutomaticEnv()
	v.Unmarshal(&cfg)

	// Validate settings
	var errb strings.Builder
	sval := reflect.ValueOf(&cfg).Elem()
	stype := sval.Type()
	for i := 0; i < sval.NumField(); i++ {
		vf := sval.Field(i)
		tf := stype.Field(i)
		// Skip any fields that don't map to an environment variable.
		envvar := tf.Tag.Get("mapstructure")
		if envvar == "" {
			continue
		}
		// Skip any fields that are not strings.
		if vf.Kind() != reflect.String {
			continue
		}
		// Skip any fields that are already set.
		if vf.Interface() != "" {
			continue
		}
		// Skip any fields that are optional.
		if ot := tf.Tag.Get("optional"); ot != "" {
			continue
		}
		// Any fields that remain, cause errors when unset.
		fmt.Fprintf(&errb, "\nenvvar %s must be set", strings.ToUpper(envvar))
	}
	if errb.Len() > 0 {
		err = errors.New(errb.String())
	}
	return
}

func (cfg config) inferior(payload []byte) (err error) {
	cmd := exec.Command(cfg.Mapper)

	// Set a signal for child processes to receive when the parent dies.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	// Create a StdinPipe and a goroutine to write to it.
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("exec stdin: %s", err)
		return
	}
	go func() {
		defer stdin.Close()
		stdin.Write(payload)
	}()

	// Wait for the child process to terminate.
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("exec wait: %s", err)
	}
	log.Printf("%s", out)
	return
}

func (cfg config) cbStan(msg *stan.Msg) {
	// Only process messages greater than the last one processed.
	// If it has been seen, skip to acknowledge it to the server.
	if msg.Sequence > cfg.LastProcessed {
		log.Printf(
			"=== {%s} processing message [seq %d] ===",
			cfg.StanChannel,
			uint64(msg.Sequence),
		)
		// Launch the mapper.
		err := cfg.inferior(msg.Data)
		if err != nil {
			return
		}
		// Upon success, record the last processed sequence number.
		atomic.SwapUint64(&cfg.LastProcessed, msg.Sequence)
	}
	msg.Ack()
}

func mainproc() (err error, usage bool) {
	// Parse CLI arguments
	cfg, err := loadcfg()
	if err != nil {
		return err, true
	}

	// Initialize NATS connection.
	nc, err := nats.Options{
		AllowReconnect: true,
		MaxReconnect:   -1,
		Name:           "nats-map",
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
		Url:            cfg.NatsUrl,
	}.Connect()
	if err != nil {
		return
	}
	defer nc.Close()

	// Initialize STAN connection.
	sc, err := stan.Connect(
		cfg.StanCluster,
		cfg.StanClientId,
		stan.NatsConn(nc),
	)
	if err != nil {
		return
	}
	defer sc.Close()

	// Main Routine.
	log.Printf(
		"Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n",
		cfg.StanChannel,
		cfg.StanClientId,
		cfg.StanQueueGroup,
		cfg.StanDurableName,
	)

	_, err = sc.QueueSubscribe(
		cfg.StanChannel,
		cfg.StanQueueGroup,
		cfg.cbStan,
		stan.StartAt(pb.StartPosition_NewOnly),
		stan.DurableName(cfg.StanDurableName),
		// Make an attempt at in-order processing.
		stan.MaxInflight(1),
		stan.SetManualAckMode(),
	)
	if err != nil {
		return
	}

	sigChan := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			log.Printf("\nInterrupted\n\n")
			done <- true
		}
	}()
	<-done
	return
}

func main() {
	log.SetFlags(0)
	ret := 0
	var err, usage = mainproc()
	if err != nil {
		log.Printf("[ERROR] %s", err)
		ret = 1
	}
	if usage {
		log.Print(usageStr)
		ret = 1
	}
	os.Exit(ret)
}
