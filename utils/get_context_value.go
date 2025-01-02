package utils

import (
	"cashflow/models"
	"context"
	"log"
)

func GetContextValue(ctx context.Context, key models.ContextKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		log.Printf("value %s not found in context", key)
		return ""
	}
	return value
}
