package connectors

import "context"

type AIConnector interface {
	GenerateResponse(ctx context.Context, message string) (string, error)
}

type OptionableConnector interface {
	AIConnector
	GenerateResponseWithOptions(ctx context.Context, message string, system string, options map[string]any) (string, error)
}

