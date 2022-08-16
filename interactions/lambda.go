package interactions

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
