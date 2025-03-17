package godatabase

import (
	"github.com/polluxdev/go-database/mongo"
	"github.com/polluxdev/go-database/mysql"
	"github.com/polluxdev/go-database/postgres"
	"github.com/polluxdev/go-database/redis"
)

type (
	PostgresConfig = postgres.Postgres
	MySQLConfig    = mysql.MySQL
	MongoConfig    = mongo.Mongo
	RedisConfig    = redis.Redis
)

var (
	NewPostgres = postgres.New
	NewMySQL    = mysql.New
	NewMongoDB  = mongo.New
	NewRedis    = redis.New
)
