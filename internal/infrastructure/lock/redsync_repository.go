package lock

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type Locker interface {
	Lock() error
	Unlock() (bool, error)
}

type redsyncRepository struct {
	mutex *redsync.Mutex
}

func NewRedsyncRepository(redisClient *redis.Client) *redsyncRepository {
	pool := goredis.NewPool(redisClient)
	redisLock := redsync.New(pool)
	mutex := redisLock.NewMutex("ratelimiter-mutex")
	return &redsyncRepository{
		mutex: mutex,
	}
}

func (r *redsyncRepository) Lock() error {
	return r.mutex.Lock()
}

func (r *redsyncRepository) Unlock() (bool, error) {
	return r.mutex.Unlock()
}
