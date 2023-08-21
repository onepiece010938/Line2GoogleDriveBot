/*
Copyright © 2023 Raymond onepiece010938@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/onepiece010938/Line2GoogleDriveBot/cmd/server"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/cache"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/ssm"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
)

var (
	ginLambda        *ginadapter.GinLambda
	cacheLambda      *cache.Cache
	lineClientLambda *linebot.Client
	ssmsvc           *ssm.SSM
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func main() {
	deploy := os.Getenv("DEPLOY_PLATFORM")
	if deploy == "lambda" {
		log.Println("deploy:", deploy)
		rootCtx, _ := context.WithCancel(context.Background()) //nolint
		ssmsvc = ssm.NewSSM()

		lineSecret, err := ssmsvc.FindParameter(rootCtx, ssmsvc.Client, "BAOSAVE_CHANNEL_SECRET")
		if err != nil {
			log.Println(err)
		}
		lineAccessToken, err := ssmsvc.FindParameter(rootCtx, ssmsvc.Client, "BAOSAVE_CHANNEL_ACCESS_TOKEN")
		if err != nil {
			log.Println(err)
		}
		lineClientLambda, err = linebot.New(lineSecret, lineAccessToken)
		if err != nil {
			log.Fatal(err)
		}

		cacheLambda = cache.NewCache(cache.InitBigCache(rootCtx))

		app := app.NewApplication(rootCtx, cacheLambda, lineClientLambda)
		ginRouter := server.InitRouter(rootCtx, app)
		ginLambda = ginadapter.New(ginRouter)

		lambda.Start(Handler)
	} else {
		log.Println("deploy: local")
		server.StartServer()
	}

}
