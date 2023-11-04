package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/dynamodb"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/google"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/ssm"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/router"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func mustGetValue(parameters map[string]string, key string) string {
	value, ok := parameters[key]
	if !ok {
		log.Fatalf("%s not found. Exiting...\n", key)
	}
	return value
}

func NewGinLambda() *ginadapter.GinLambda {
	rootCtx, _ := context.WithCancel(context.Background()) //nolint
	ssmsvc := ssm.NewSSM()

	parameterNames := []string{
		"BAOSAVE_CHANNEL_SECRET",
		"BAOSAVE_CHANNEL_ACCESS_TOKEN",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"REDIRECT_URL",
	}

	parameters, err := ssmsvc.FindParameters(rootCtx, ssmsvc.Client, parameterNames)
	if err != nil {
		log.Println(err)
	}

	lineSecret := mustGetValue(parameters, "BAOSAVE_CHANNEL_SECRET")
	lineAccessToken := mustGetValue(parameters, "BAOSAVE_CHANNEL_ACCESS_TOKEN")
	googleClientID := mustGetValue(parameters, "GOOGLE_CLIENT_ID")
	googleClientSecret := mustGetValue(parameters, "GOOGLE_CLIENT_SECRET")
	redirectURL := mustGetValue(parameters, "REDIRECT_URL")

	lineClientLambda, err := linebot.New(lineSecret, lineAccessToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("LineBot Create Success")

	db := dynamodb.NewTableBasics("google-oauth")

	oauth := google.NewGoogleOAuth(googleClientID, googleClientSecret, redirectURL)

	app := app.NewApplication(rootCtx, db, oauth, lineClientLambda)
	ginRouter := initRouter(rootCtx, app)
	return ginadapter.New(ginRouter)
}

func initRouter(rootCtx context.Context, app *app.Application) (ginRouter *gin.Engine) {

	// Create gin router
	ginRouter = gin.New()

	// Set general middleware
	// router.SetGeneralMiddlewares(rootCtx, ginRouter)

	// Register all handlers
	router.RegisterHandlers(ginRouter, app)

	return ginRouter
}

func StartNgrokServer() {
	rootCtx, rootCtxCancelFunc := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	lineClient, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err.Error())
	}

	db := dynamodb.NewTableBasics("google-oauth")
	// change to local dynamodb
	db.DynamoDbClient = dynamodb.CreateLocalClient(8000)

	oauth := google.NewGoogleOAuth(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("REDIRECT_URL"))

	app := app.NewApplication(rootCtx, db, oauth, lineClient)

	ginRouter := initRouter(rootCtx, app)
	// Run ngrok
	wg.Add(1)
	runNgrokServer(rootCtx, &wg, ginRouter)

	// Listen to SIGTERM/SIGINT to close
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	<-gracefulStop
	rootCtxCancelFunc()

	// Wait for all services to close with a specific timeout
	var waitUntilDone = make(chan struct{})
	go func() {
		wg.Wait()
		close(waitUntilDone)
	}()
	select {
	case <-waitUntilDone:
		log.Println("success to close all services")
	case <-time.After(10 * time.Second):
		log.Println(context.DeadlineExceeded, "fail to close all services")
	}

}
func runNgrokServer(rootCtx context.Context, wg *sync.WaitGroup, ginRouter *gin.Engine) {

	tun, err := ngrok.Listen(rootCtx,
		config.HTTPEndpoint(config.WithDomain(os.Getenv("NGROK_DOMAIN"))),
		ngrok.WithAuthtokenFromEnv(),
	)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Application available at:", tun.URL())

	// Run the server in a goroutine
	go func() {
		err = http.Serve(tun, ginRouter)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for rootCtx done
	go func() {
		<-rootCtx.Done()
		// Create a context with a timeout for closing the ngrok tunnel
		log.Println("Shutting down ngrok server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Use the created context to close the ngrok tunnel with a timeout
		if err := tun.CloseWithContext(ctx); err != nil {
			log.Printf("Error closing ngrok tunnel: %v\n", err)
		}
		log.Println("ngrok server gracefully stopped")
		wg.Done()
	}()

}
