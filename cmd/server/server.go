package main

import (
	"database/sql"
	"fmt"
	"github.com/carlosroman/payments-api/internal/app/payment"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"time"
)

const LOCAL_BUILD_VERSION = "snapshot"

var version = LOCAL_BUILD_VERSION

func main() {
	app := cli.NewApp()
	app.Name = "Payment server"
	app.Version = version
	app.Authors = []cli.Author{
		{
			Name:  "Carlos Roman",
			Email: "carlosr@cliche-corp.co.uk",
		},
	}
	log.SetLevel(log.InfoLevel)

	app.Commands = []cli.Command{
		{Name: "run",
			Aliases: []string{"r"},
			Usage:   "run server",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "port, p",
					Value:  8080,
					Usage:  "Set the port of the server",
					EnvVar: "SERVER_PORT",
				},
				cli.StringFlag{
					Name:   "db-user",
					Usage:  "Database username",
					EnvVar: "DB_USER",
				},
				cli.StringFlag{
					Name:   "db-password",
					Usage:  "Database password",
					EnvVar: "DB_PASSWORD",
				},
				cli.StringFlag{
					Name:   "db-name",
					Usage:  "Database name",
					EnvVar: "DB_NAME",
				},
				cli.StringFlag{
					Name:   "db-host",
					Value:  "localhost",
					Usage:  "Database host",
					EnvVar: "DB_HOST",
				},
				cli.IntFlag{
					Name:   "db-port",
					Value:  5432,
					Usage:  "The database port",
					EnvVar: "DB_PORT",
				},
			},
			Action: func(c *cli.Context) error {
				db, err := initDb(
					c.String("db-host"),
					c.Int("db-port"),
					c.String("db-user"),
					c.String("db-password"),
					c.String("db-name"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				dir, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				log.Infof("current dir: %s", dir)

				s := payment.NewService(db)
				h := payment.GetHandlers(s)

				h.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
				addr := fmt.Sprintf("0.0.0.0:%v", c.Int("port"))
				headersOk := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
				originsOk := handlers.AllowedOrigins([]string{"*"})
				methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
				srv := &http.Server{
					Addr: addr,
					// Good practice to set timeouts to avoid Slowloris attacks.
					WriteTimeout: time.Second * 15,
					ReadTimeout:  time.Second * 15,
					IdleTimeout:  time.Second * 60,
					Handler:      handlers.CORS(headersOk, originsOk, methodsOk)(h), // Pass our instance of gorilla/mux in.
				}

				log.Infof("Starting server at %s", addr)
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

//func corsHandler(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//		h.ServeHTTP(w, r)
//	})
//}

func initDb(dbHost string, dbPort int, dbUser string, dbPassword string, dbName string) (db *sql.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return db, err
	}
	log.Info("Successfully connected!")
	return db, err
}
