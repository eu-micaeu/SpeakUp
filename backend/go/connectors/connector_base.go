package connectors

import "context"

type AIConnector interface {
	GenerateResponse(ctx context.Context, prompt string) (string, error)
}
