package redis


import (
  "context"
  "flag"
  "fmt"
  "time"


  "github.com/hoangminhphuc/goph-chat/common/logger"
  "github.com/redis/go-redis/v9"
)


var (
  defaultRedisDB     = 0
  defaultRedisMaxActive = 0 // 0 is unlimited max active connection
  defaultRedisMaxIdle   = 10
  defaultDialTimeout = 5 * time.Second
  defaultReadTimeout = 3 * time.Second
  defaultWriteTimeout = 3 * time.Second
  defaultPoolTimeout = 4 * time.Second
)




type RedisOption struct {
  addr        string
  db          int
  dialTimeout time.Duration // Timeout for initial TCP connection
  readTimeout time.Duration // Timeout for reading a single command reply
  writeTimeout time.Duration // Timeout for writing a single command
  poolSize    int
  maxIdleConns int // Maximum number of idle connections in the pool      
  poolTimeout time.Duration // Time a client waits for a free connection before timeout
}


type RedisService struct {
    name  string
    client      *redis.Client
    RedisOption
    logger logger.ZapLogger
}


func NewRedisDB() *RedisService {
  return &RedisService{
      name: "redis",
      RedisOption: RedisOption{
        db : defaultRedisDB,
        poolSize: defaultRedisMaxActive,
        maxIdleConns: defaultRedisMaxIdle,
      },
      logger: logger.NewZapLogger(),
  }
}


func (r *RedisService) Name() string {
  return r.name
}


func (r *RedisService) InitFlags() {
  prefix := r.name
  if r.name != "" {
      prefix += "-"
  }
    flag.StringVar(&r.addr, prefix+"addr",
        "localhost:6379", "Redis address (host:port)")

    flag.IntVar(&r.db, prefix+"db",
        defaultRedisDB, "Redis DB number")

    flag.DurationVar(&r.dialTimeout, prefix+"dial-timeout",
        defaultDialTimeout, "Dial timeout")

    flag.DurationVar(&r.readTimeout, prefix+"read-timeout",
        defaultReadTimeout, "Read timeout")

    flag.DurationVar(&r.writeTimeout, prefix+"write-timeout",
        defaultWriteTimeout, "Write timeout")

    flag.IntVar(&r.poolSize, prefix+"pool-size",
        defaultRedisMaxActive, "Connection pool size")

    flag.IntVar(&r.maxIdleConns, prefix+"max-idle-conns",
        defaultRedisMaxIdle, " Maximum number of idle connections")

    flag.DurationVar(&r.poolTimeout, prefix+"pool-timeout",
        defaultPoolTimeout, "Timeout for getting a connection from the pool")
}


// Run initializes and verifies the Redis client
func (r *RedisService) Run() error {
    r.logger.Log.Info("Connecting to Redis at ", r.addr, "...")

    opt := &redis.Options{
			Addr:         r.addr,
			DB:           r.db,
			DialTimeout:  r.dialTimeout,
			ReadTimeout:  r.readTimeout,
			WriteTimeout: r.writeTimeout,
			PoolSize:     r.poolSize,
			MinIdleConns: r.maxIdleConns,
			PoolTimeout:  r.poolTimeout,
	}

		client := redis.NewClient(opt)

    if err := client.Ping(context.Background()).Err(); err != nil {
			r.logger.Log.Error("Cannot connect to Redis:", err.Error())
			return fmt.Errorf("failed to connect to Redis: %w", err)
		}

		r.client = client
		r.logger.Log.Info("Redis connection established successfully")

    return nil
}

func (r *RedisService) Stop() <-chan error {
  ch := make(chan error, 1)

  go func() {
    var err error
    if r.client != nil {
        if cerr := r.client.Close(); cerr != nil {
            r.logger.Log.Info("cannot close ", r.name, " error:", cerr)
            err = cerr
        }
    }
    ch <- err
    close(ch)
  }()

	return ch
}


func (r *RedisService) Get() interface{} {
  return r.client
}


