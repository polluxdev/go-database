package dbx

import (
	"github.com/polluxdev/go-dbx/mongo"
	"github.com/polluxdev/go-dbx/mysql"
	"github.com/polluxdev/go-dbx/postgres"
	"github.com/polluxdev/go-dbx/redis"
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
