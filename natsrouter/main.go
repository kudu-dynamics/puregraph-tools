/*
 * natsrouter
 * ==========
 *
 * Webhook to reroute incoming requests to NATS.
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

type ApiCtx struct {
	prefix string
	nc     *nats.Conn
}

func (a ApiCtx) Serve(c echo.Context) (err error) {
	// Decode JSON input.
	m := echo.Map{}
	if err = c.Bind(&m); err != nil {
		return err
	}

	// Do not reroute empty messages or messages without a routing key.
	if key, ok := m["Key"].(string); len(m) > 0 && ok {
		// Re-encode the request body to bytes.
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}

		// Route by the given key.
		// e.g.
		// prefix  = "jupyter"
		// key     = "data-bysha256/{...}"
		// channel = "jupyter.data-bysha256"
		keyparts := strings.Split(key, "/")
		channel := fmt.Sprintf("%s.%s", a.prefix, keyparts[0])

		// Publish a message to the NATS channel.
		// DEV: If the NATS Connection is reconnecting, the client will store some
		//      number of messages in an in-memory queue until full while
		//      indefinitely retrying to connect. An error will only materialize
		//      when the internal buffer is full.
		if err = a.nc.Publish(channel, b); err != nil {
			return err
		}
	}

	// DEV: MinIO expects a body for a 200 response.
	//      Otherwise connections will remain ESTABLISHED indefinitely.
	//
	//      https://github.com/minio/minio/issues/8435
	return c.String(http.StatusOK, "ok")
}

func main() {
	// Parse CLI arguments
	v := viper.New()
	v.SetDefault("nats_url", "nats://localhost:4222")
	v.SetDefault("nats_prefix", "test")
	v.SetDefault("port", "9090")
	v.AutomaticEnv()

	// Initialize new components
	e := echo.New()

	nc, err := nats.Options{
		AllowReconnect: true,
		MaxReconnect:   -1,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
		Url:            v.GetString("nats_url"),
	}.Connect()

	if err != nil {
		e.Logger.Fatal(err)
	}

	ac := &ApiCtx{prefix: v.GetString("nats_prefix"), nc: nc}

	// Add Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Add Handlers
	e.POST("/", ac.Serve)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", v.GetString("port"))))
}
