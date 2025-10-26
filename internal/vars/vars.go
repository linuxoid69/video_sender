package vars

type Config struct {
	TelegramToken string `env:"VS_TG_TOKEN,required"`
	TelegramGroup int64  `env:"VS_TG_GROUP,required"`
	RedisHost     string `env:"VS_REDIS_HOST,required"`
	RedisPassword string `env:"VS_REDIS_PASSWORD"`
}
