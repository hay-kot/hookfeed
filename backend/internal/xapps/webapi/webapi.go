package webapi

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/web/docs"
	"github.com/hay-kot/hookfeed/backend/internal/web/mid"
	"github.com/hay-kot/hookfeed/backend/internal/xapps/webapi/handlers"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Config struct {
	Host           string        `toml:"host"            env:"WEB_HOST"            envDefault:"0.0.0.0"`
	Port           string        `toml:"port"            env:"WEB_PORT"            envDefault:"9990"`
	AllowedOrigins []string      `toml:"allowed_origins" env:"WEB_ALLOWED_ORIGINS" envDefault:"*"`
	TLSCert        string        `toml:"tls_cert"        env:"WEB_TLS_CERT"        envDefault:""`
	TLSKey         string        `toml:"tls_key"         env:"WEB_TLS_KEY"         envDefault:""`
	IdleTimeout    time.Duration `toml:"idle_timeout"    env:"WEB_IDLE_TIMEOUT"    envDefault:"30s"`
	ReadTimeout    time.Duration `toml:"read_timeout"    env:"WEB_READ_TIMEOUT"    envDefault:"10s"`
	WriteTimeout   time.Duration `toml:"write_timeout"   env:"WEB_WRITE_TIMEOUT"   envDefault:"20s"`
}

func (cfg Config) Addr() string {
	return cfg.Host + ":" + cfg.Port
}

type WebAPI struct {
	l        zerolog.Logger
	build    string
	cfg      Config
	services *services.Service
}

func New(l zerolog.Logger, build string, cfg Config, services *services.Service) *WebAPI {
	return &WebAPI{
		l:        l.With().Str("service", "web_api").Logger(),
		cfg:      cfg,
		build:    build,
		services: services,
	}
}

func (ib *WebAPI) Start(ctx context.Context) error {
	mux := ib.routes()

	server := &http.Server{
		Handler:      mux,
		Addr:         ib.cfg.Addr(),
		IdleTimeout:  ib.cfg.IdleTimeout,
		ReadTimeout:  ib.cfg.ReadTimeout,
		WriteTimeout: ib.cfg.WriteTimeout,
	}

	go func() {
		<-ctx.Done()
		ib.l.Info().Msg("stopping service")
		_ = server.Shutdown(context.Background())
	}()

	ib.l.Info().Str("docs", "http://"+ib.cfg.Addr()+"/docs/index.html").Msg("starting service")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (ib *WebAPI) routes() chi.Router {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Use(
		middleware.RealIP,
		middleware.CleanPath,
		middleware.StripSlashes,
		mid.RequestID(),
		mid.Logger(ib.l, ""),
		middleware.AllowContentType("application/json", "text/plain", "text/html"),
	)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   ib.cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Trace-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	adapter := mid.ErrorHandler(ib.l)

	userctrl := handlers.NewAuthController(ib.services.Users, ib.services.Passwords)
	feedctrl := handlers.NewFeedController(ib.services.Feeds)

	mux.HandleFunc("GET /docs/swagger.json", adapter.Adapt(docs.SwaggerJSON))
	mux.HandleFunc("GET /api/v1/info", adapter.Adapt(handlers.Info(dtos.StatusResponse{Build: ib.build})))

	mux.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.PersistAuthorization(true),
		httpSwagger.URL("/docs/swagger.json"),
	))

	mux.Post("/api/v1/users/login", adapter.Adapt(userctrl.Authenticate))
	mux.Post("/api/v1/users/register", adapter.Adapt(userctrl.Register))
	mux.Post("/api/v1/users/reset-password-request", adapter.Adapt(userctrl.ResetPasswordRequest))
	mux.Post("/api/v1/users/reset-password", adapter.Adapt(userctrl.ResetPassword))

	mux.Group(func(r chi.Router) {
		r.Use(mid.Authenticate(ib.services.Users))

		r.Get("/api/v1/users/self", adapter.Adapt(userctrl.Self))
		r.Patch("/api/v1/users/self", adapter.Adapt(userctrl.Update))

		r.Get("/api/v1/feeds", adapter.Adapt(feedctrl.GetAll))

		// $scaffold_inject_routes
	})

	return mux
}
