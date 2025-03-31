package connectors

import "context"

type AIConnector interface {
	GenerateResponse(ctx context.Context, message string) (string, error)
}
