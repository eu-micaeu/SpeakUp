package planlimits

import (
	"context"
)

func GetFreeDailyLimit() int64 {
	return -1
}

func IsProUser(ctx context.Context, userID string) (bool, error) {
	return true, nil
}

func GetUsageCount(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func CheckAndIncrementUsage(ctx context.Context, userID string, limit int64) (bool, error) {
	return true, nil
}

