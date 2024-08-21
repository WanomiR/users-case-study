package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wanomir/e"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	// include to use db drivers
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var appInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "geoservice",
	Name:      "info",
	Help:      "App environment info",
}, []string{"version"})

type Config struct {
	host       string
	port       string
	dsn        string
	appVersion string
}

type App struct {
	config     Config
	server     *http.Server
	signalChan chan os.Signal
	//DB          repository.Repository
}

func NewApp() (a *App, err error) {
	defer func() { err = e.WrapIfErr("failed to init app", err) }()

	a = &App{}

	if err = a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start() {
	//defer a.DB.Connection().Close()

	fmt.Println("Started server on port", a.config.port)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (a *App) Shutdown() {
	<-a.signalChan

	// a-la graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	<-ctx.Done()

	fmt.Println("Shutting down server gracefully")
}

func (a *App) readConfig() (err error) {

	a.config.host = os.Getenv("HOST")
	a.config.port = os.Getenv("PORT")

	a.config.dsn = fmt.Sprintf( // database source name
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5\n",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB_NAME"),
	)

	a.config.appVersion = os.Getenv("APP_VERSION")

	if a.config.host == "" || a.config.port == "" || a.config.dsn == "" || a.config.appVersion == "" {
		return errors.New("env variables not set")
	}

	return nil
}

func (a *App) init() error {
	if err := a.readConfig(); err != nil {
		return err
	}

	_, err := a.connectToDB()
	if err != nil {
		return err
	}

	a.server = &http.Server{
		Addr:         ":" + a.config.port,
		Handler:      a.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

	appInfo.With(prometheus.Labels{"version": a.config.appVersion}).Set(1)

	return nil
}

func (a *App) connectToDB() (conn *sql.DB, err error) {
	defer func() { err = e.WrapIfErr("failed to connect to database", err) }()

	conn, err = sql.Open("pgx", a.config.dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
