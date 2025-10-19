package queue

import "context"

type Queuer interface {
	CreateJob(ctx context.Context, key string, value any) error
	GetJob(ctx context.Context, key string) (string, error)
	DeleteJob(ctx context.Context, keys ...string) error
}
