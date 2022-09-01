package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"pairprofit.com/x/handlers/auth"
	"pairprofit.com/x/handlers/communication"
	"pairprofit.com/x/handlers/listings"
	"pairprofit.com/x/helpers"
	"pairprofit.com/x/middleware"
)

var Version = "unknown"

func run() error {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	if gin.IsDebugging() {
		log.Print("Using debugging mode " + Version)
		if err := godotenv.Load("app.env"); err != nil {
			return err
		}
	} else {
		log.Print("Using release mode " + Version)

		sentryDsn := helpers.GetEnv("SENTRY_DSN", "")
		if len(sentryDsn) > 0 {
			log.Print("Starting Sentry")
		}
	}

	os.Setenv("AWS_PROFILE", "personal")

	r.Use(
		middleware.ConfigMiddleware(),
		middleware.Sesv2Middleware(),
		middleware.S3Middleware(),
		middleware.CORSMiddleware(),
		middleware.RedisMiddleware(),
		middleware.CognitoMiddleware(Version),
		middleware.SIBMiddleware(),
	)

	// Auth Endpoints
	r.POST("/auth/login", auth.Login)
	r.POST("/auth/link", auth.LoginViaLink)
	r.GET("/auth/logout", auth.TokenRevoke)
	r.POST("/auth/refresh", auth.TokenRefresh)
	r.POST("/auth/register", auth.SignUp)
	r.POST("/auth/forgot_password", auth.ForgotPassword)
	r.POST("/auth/reset_password", auth.ResetPassword)
	r.POST("/auth/update_password", auth.UpdatePassword)
	r.GET("/auth/access", auth.VerifyAuth)

	r.GET("/app/email", auth.EmailEndPoint)

	r.GET("/pb", listings.GetListings)

	//websocket
	r.GET("/ws/:id", communication.WebsocketHandler)

	// Ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// build and start server
	log.Print("Listening at :8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Panicln("Server Shutdown: ...", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("Server exiting ...")
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
