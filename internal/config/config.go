package config

type Variables struct {
	RedisPriceAddress          string `env:"REDIS_PRICE_ADDRESS"`
	RedisPricePassword         string `env:"REDIS_PRICE_PASSWORD"`
}