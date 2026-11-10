package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	authApp "github.com/hguitar/backend-go/internal/auth/application"
	authInfra "github.com/hguitar/backend-go/internal/auth/infrastructure/persistence"
	authSecurity "github.com/hguitar/backend-go/internal/auth/infrastructure/security"
	authPres "github.com/hguitar/backend-go/internal/auth/presentation"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/config"
	contentApp "github.com/hguitar/backend-go/internal/content/application"
	contentInfra "github.com/hguitar/backend-go/internal/content/infrastructure/persistence"
	contentPres "github.com/hguitar/backend-go/internal/content/presentation"
	learningApp "github.com/hguitar/backend-go/internal/learning/application"
	learningInfra "github.com/hguitar/backend-go/internal/learning/infrastructure/persistence"
	learningPres "github.com/hguitar/backend-go/internal/learning/presentation"
	musicApp "github.com/hguitar/backend-go/internal/music/application"
	musicInfra "github.com/hguitar/backend-go/internal/music/infrastructure/persistence"
	musicPres "github.com/hguitar/backend-go/internal/music/presentation"
	onboardingApp "github.com/hguitar/backend-go/internal/onboarding/application"
	onboardingInfra "github.com/hguitar/backend-go/internal/onboarding/infrastructure/persistence"
	onboardingPres "github.com/hguitar/backend-go/internal/onboarding/presentation"
	progressApp "github.com/hguitar/backend-go/internal/progress/application"
	progressInfra "github.com/hguitar/backend-go/internal/progress/infrastructure/persistence"
	progressPres "github.com/hguitar/backend-go/internal/progress/presentation"
	songApp "github.com/hguitar/backend-go/internal/song/application"
	songInfra "github.com/hguitar/backend-go/internal/song/infrastructure/persistence"
	songPres "github.com/hguitar/backend-go/internal/song/presentation"
	userApp "github.com/hguitar/backend-go/internal/user/application"
	userInfra "github.com/hguitar/backend-go/internal/user/infrastructure/persistence"
	userPres "github.com/hguitar/backend-go/internal/user/presentation"
	_ "github.com/denisenkom/go-mssqldb"
)

type passwordHasherAdapter struct{}

func (passwordHasherAdapter) Hash(p string) string  { return authSecurity.HashPassword(p) }
func (passwordHasherAdapter) Check(p, h string) bool { return authSecurity.CheckPassword(p, h) }

func main() {
	cfg := config.Load()
	slog.Info("starting backend-go", "port", cfg.ServerPort)

	db, err := sql.Open("sqlserver", cfg.DBURL)
	if err != nil {
		slog.Error("failed to open db", "err", err)
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		slog.Error("failed to ping db", "err", err)
		panic(err)
	}
	slog.Info("database connected")

	// ── Auth ──
	userRepo := authInfra.NewUserRepositoryAdapter(db)
	refreshTokenRepo := authInfra.NewRefreshTokenRepositoryAdapter(db)
	jwtCfg := authSecurity.JWTConfig{
		Secret:       cfg.JWTSecret,
		AccessExpiry: cfg.AccessExpiry,
		RefreshExpiry: cfg.RefreshExpiry,
	}
	userProfileRepo := userInfra.NewUserProfileRepositoryAdapter(db)
	userProfileSvc := userApp.NewUserProfileService(userProfileRepo)

	authSvc := authApp.NewAuthService(
		userRepo, refreshTokenRepo, passwordHasherAdapter{}, jwtCfg, userProfileSvc, cfg.RefreshExpiry.Milliseconds(),
	)
	firebaseVerifier := authSecurity.NewFirebaseTokenVerifier(cfg.FirebaseProjectID)
	authHandler := authPres.NewAuthHandler(authSvc, firebaseVerifier)

	// ── User ──
	userHandler := userPres.NewUserHandler(userProfileSvc)

	// ── Content ──
	contentRepo := contentInfra.NewLessonBlockRepositoryAdapter(db)
	contentSvc := contentApp.NewContentService(contentRepo)
	contentHandler := contentPres.NewContentHandler(contentSvc)

	// ── Music ──
	chordRepo := musicInfra.NewChordRepositoryAdapter(db)
	chordSvc := musicApp.NewChordService(chordRepo)
	strummingRepo := musicInfra.NewStrummingRepositoryAdapter(db)
	strummingInstrRepo := musicInfra.NewStrummingInstructionRepositoryAdapter(db)
	strummingSvc := musicApp.NewStrummingPatternService(strummingRepo, strummingInstrRepo)
	musicHandler := musicPres.NewMusicHandler(chordSvc, strummingSvc)

	// ── Learning ──
	courseRepo := learningInfra.NewCourseRepositoryAdapter(db)
	moduleRepo := learningInfra.NewModuleRepositoryAdapter(db)
	lessonRepo := learningInfra.NewLessonRepositoryAdapter(db)
	learningMethodRepo := learningInfra.NewLearningMethodRepositoryAdapter(db)
	learningQuerySvc := learningApp.NewLearningQueryService(courseRepo, moduleRepo, lessonRepo)
	learningMethodSvc := learningApp.NewLearningMethodService(learningMethodRepo)
	learningHandler := learningPres.NewLearningHandler(learningQuerySvc, learningMethodSvc)

	// ── Progress ──
	progressRepo := progressInfra.NewLessonProgressRepositoryAdapter(db)
	progressSvc := progressApp.NewProgressService(progressRepo, lessonRepo)
	progressHandler := progressPres.NewProgressHandler(progressSvc)

	// ── Song ──
	songRepo := songInfra.NewSongRepositoryAdapter(db)
	songSvc := songApp.NewSongService(songRepo)
	songHandler := songPres.NewSongHandler(songSvc)

	// ── Onboarding ──
	onboardingRepo := onboardingInfra.NewUserPreferenceRepositoryAdapter(db)
	onboardingSvc := onboardingApp.NewOnboardingService(onboardingRepo)
	onboardingHandler := onboardingPres.NewOnboardingHandler(onboardingSvc)

	// ── Router ──
	r := chi.NewRouter()
	r.Use(common.Recovery)
	r.Use(common.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	verifyFunc := func(tokenStr string) (string, string, []string, error) {
		claims, err := authSecurity.ValidateJWT(cfg.JWTSecret, tokenStr)
		if err != nil {
			return "", "", nil, err
		}
		return claims.UserID, claims.Email, claims.Roles, nil
	}
	r.Use(common.AuthMiddleware(verifyFunc))

	api := chi.NewRouter()

	api.Mount("/auth", authHandler.Routes())
	api.Mount("/users", userHandler.Routes())
	api.Route("/users/me/progress", func(r chi.Router) {
		r.Get("/", progressHandler.GetProgress)
	})
	api.Route("/lessons/{id}/complete", func(r chi.Router) {
		r.Post("/", progressHandler.CompleteLesson)
	})
	api.Mount("/lessons", contentHandler.Routes())
	api.Mount("/chords", musicHandler.ChordRoutes())
	api.Mount("/strumming-patterns", musicHandler.StrummingRoutes())
	api.Mount("/", learningHandler.Routes())
	api.Mount("/songs", songHandler.Routes())
	api.Mount("/onboarding", onboardingHandler.Routes())

	r.Mount("/api/v1", api)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	slog.Info("server listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "err", err)
		panic(err)
	}
}
