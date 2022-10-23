package connectors

import "github.com/go-redis/redis/v8"

type RedisConnector struct {
	client *redis.Client
}

func (r *RedisConnector) Init(host string, password string, db int) {
	r.client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
}
