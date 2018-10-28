package main

import (
	"fmt"
	"github.com/carlosroman/payments-api/internal/app/payment"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"time"
)

func main() {
	var port int
	app := cli.NewApp()
	app.Name = "Payment server"

	app.Commands = []cli.Command{
		{Name: "run",
			Aliases: []string{"r"},
			Usage:   "run server",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "port, p",
					Value:       8080,
					Usage:       "Set the port of the server",
					EnvVar:      "PORT",
					Destination: &port,
				}},
			Action: func(c *cli.Context) error {
				s := payment.NewService()
				h := payment.GetHandlers(s)
				srv := &http.Server{
					Addr: fmt.Sprintf("0.0.0.0:%v", port),
					// Good practice to set timeouts to avoid Slowloris attacks.
					WriteTimeout: time.Second * 15,
					ReadTimeout:  time.Second * 15,
					IdleTimeout:  time.Second * 60,
					Handler:      h, // Pass our instance of gorilla/mux in.
				}

				if err := srv.ListenAndServe(); err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
