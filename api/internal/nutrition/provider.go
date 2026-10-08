package nutrition

import "context"

type Provider interface {
	GenerateNutrition(context.Context, Request) (*GeneratedFood, error)
}
