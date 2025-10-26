package queue

import "context"

type Queuer interface {
	CreateJob(ctx context.Context, key string, value any) error
	GetJob(ctx context.Context, key string) (string, error)
	DeleteJob(ctx context.Context, keys ...string) error
	GetKeys(ctx context.Context, pattern string) ([]string, error)
}

type Query struct {
	Key   string    `json:"key"`
	TTL   int       `json:"ttl,omitempty"`
	Value VideoData `json:"value"`
}

type VideoData struct {
	FileSize   int64  `json:"file_size"`
	VideoFile  string `json:"video_file"`
	CameraName string `json:"camera_name"`
}
