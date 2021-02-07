package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ersmith/mailgun-coding-challenge/config"
	"github.com/ersmith/mailgun-coding-challenge/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

const JSON = "application/json"

type App struct {
	Router   *mux.Router
	dbPool   *pgxpool.Pool
	validate *validator.Validate
	Logger   *zap.SugaredLogger
}

// Initializes the application and sets up items like the
// database connections
func (self *App) Initialize(dbConfig *config.DbConfig, logger *zap.SugaredLogger) {
	self.Logger = logger
	self.initializeDb(dbConfig)
	self.initializeRouter()
	self.validate = validator.New()
}

// Starts the app. This will end with it listening on the specified port
func (self *App) Run(port string) {
	self.Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), self.Router))
}

// Performs any cleanup needed for the app
func (self *App) Cleanup() {
	self.dbPool.Close()
}

// Performs the database initialization
func (self *App) initializeDb(dbConfig *config.DbConfig) {
	dbUrl := dbConfig.ConnnectionUrl()
	self.Logger.Infow("Connecting to database",
		zap.String("user", dbConfig.Username),
		zap.String("host", dbConfig.Host),
		zap.String("port", dbConfig.Port),
		zap.String("name", dbConfig.Name))
	config, err := pgxpool.ParseConfig(dbUrl)

	if err != nil {
		self.Logger.Fatalw("unable to parse database config", zap.Error(err))
	}

	config.MinConns = dbConfig.MinPoolSize
	config.MaxConns = dbConfig.MaxPoolSize
	config.ConnConfig.Logger = zapadapter.NewLogger(self.Logger.Desugar())

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		self.Logger.Fatalw("unable to connect to database", zap.Error(err))
		os.Exit(1)
	}

	self.dbPool = pool
}

// Instantiates the mux router and configures the routing
func (self *App) initializeRouter() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(self.loggingMiddleware)
	eventsSubrouter := router.PathPrefix("/events/").Subrouter()
	eventsSubrouter.Use(self.domainValidationMiddleware)
	eventsSubrouter.HandleFunc("/{domain}/delivered", self.putEventDeliveredHandler).Methods("PUT")
	eventsSubrouter.HandleFunc("/{domain}/bounced", self.putEventBouncedHandler).Methods("PUT")

	domainsSubrouter := router.PathPrefix("/domains/").Subrouter()
	domainsSubrouter.Use(self.domainValidationMiddleware)
	domainsSubrouter.HandleFunc("/{domain}", self.getDomainHandler).Methods("GET")

	self.Router = router
}

// Handles requests for delivered events
func (self *App) putEventDeliveredHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	err := self.validate.Var(domainName, "fqdn")

	if err != nil {
		self.Logger.Infow("bad domain name", zap.Error(err), zap.String("domain", domainName))
		badRequest(w)
		return
	}

	self.Logger.Infow("Delivered event", zap.String("domain", domainName))
	domain := models.Domain{}
	domain.DomainName = domainName
	err = domain.IncrementDelivered(self.dbPool)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Handles requests for dropped events
func (self *App) putEventBouncedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	self.Logger.Infow("bounced event", zap.String("domain", vars["domain"]))
	domain := models.Domain{}
	domain.DomainName = vars["domain"]
	err := domain.IncrementBounced(self.dbPool)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Handles requests for info about a domain
func (self *App) getDomainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	self.Logger.Infow("get domain", zap.String("domain", vars["domain"]))
	domain, err := models.GetDomain(self.dbPool, self.Logger, vars["domain"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", JSON)
	json.NewEncoder(w).Encode(domain.Json())
}

// Logs all incoming requests
func (self *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		self.Logger.Infow("received request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method))

		next.ServeHTTP(w, r)
	})
}

// Middleware for validating the domain in the path
func (self *App) domainValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainName := mux.Vars(r)["domain"]
		err := self.validate.Var(domainName, "fqdn")

		if err != nil {
			self.Logger.Infow("bad domain name", zap.Error(err), zap.String("domain", domainName))
			badRequest(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Sets up a bad request response
func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}
