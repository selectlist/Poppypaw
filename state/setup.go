package state

import (
	"errors"
	"fmt"
	"poppypaw/config"

	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/infinitybotlist/eureka/genconfig"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

func Setup() {
	fmt.Println("Generating config data")
	genconfig.GenConfig(config.Config{})

	fmt.Println("Loading: config.yaml")
	cfg, err := os.ReadFile("config.yaml")

	if err != nil {
		fmt.Println("Uh oh! Failed to read config:", err.Error())
	}

	err = yaml.Unmarshal(cfg, &Config)

	if err != nil {
		fmt.Errorf("Uh oh! Failed to unmarshal config:", err.Error())
	}

	fmt.Println("Validating config data")
	err = Validator.Struct(Config)

	if err != nil {
		fmt.Println("Uh oh! The config validation failed:", err.Error())
	}

	fmt.Println("Connecting to Discord")
	Discord, err = createDiscordSession()

	if err != nil {
		fmt.Println("Error creating Discord session:", err)
	}

	fmt.Println("Connecting to Database")
	Database, err = createDatabase()

	if err != nil {
		fmt.Errorf("Error creating Database pool:", err)
	}

	fmt.Println("Connecting to Redis")
	Redis, err = createRedisClient()

	if err != nil {
		fmt.Errorf("Error creating Redis client:", err)
	}
}

func createDiscordSession() (*discordgo.Session, error) {
	dg, err := discordgo.New("Bot " + Config.DiscordToken)
	if err != nil {
		return nil, err
	}

	return dg, nil
}

func createDatabase() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(Context, Config.DatabaseURL)
	if err != nil {
		return nil, errors.New("Uh oh! Unable to connect to Database: " + err.Error())
	}

	return pool, nil
}

func createRedisClient() (*redis.Client, error) {
	opt, err := redis.ParseURL(Config.RedisURL)
	if err != nil {
		return nil, errors.New("Uh oh! Failed to connect to Redis: " + err.Error())
	}
	client := redis.NewClient(opt)

	return client, nil
}
