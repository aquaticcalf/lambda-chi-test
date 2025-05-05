package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
)

// chiLambda is a global variable to store the chi router instance
// it is used to process API Gateway proxy requests through the chi router.
var chiLambda *chiadapter.ChiLambda

// Handler processes API Gateway requests through the Chi router.
// Initializes the router on first invocation (cold start).
func Handler(router *chi.Mux) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if chiLambda == nil {
		log.Printf("Cold start: initializing Chi router")
		chiLambda = chiadapter.New(router)
	}

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return chiLambda.ProxyWithContext(ctx, req)
	}
}

// Accepts a Chi router directly for type safety.
func Start(router *chi.Mux) {
	log.Printf("Starting in Lambda mode")
	lambda.Start(Handler(router))
}
