package state

import (
	"errors"
	"fmt"
	"poppypaw/config"

	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/infinitybotlist/eureka/genconfig"
	"github.com/jackc/pgx/v5/pgxpool"
	novu "github.com/novuhq/go-novu/lib"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

func Setup() {
	// Config
	fmt.Println("Generating config data")
	genconfig.GenConfig(config.Config{})

	fmt.Println("Loading: config.yaml")
	cfg, err := os.ReadFile("config.yaml")

	if err != nil {
		sentry.CaptureException(err)
	}

	err = yaml.Unmarshal(cfg, &Config)

	if err != nil {
		sentry.CaptureException(err)
	}

	fmt.Println("Validating config data")
	err = Validator.Struct(Config)

	if err != nil {
		sentry.CaptureException(err)
	}

	// Discord
	fmt.Println("Connecting to Discord")
	Discord, err = createDiscordSession()

	if err != nil {
		sentry.CaptureException(err)
	}

	// Database
	fmt.Println("Connecting to Database")
	Database, err = createDatabase()

	if err != nil {
		sentry.CaptureException(err)
	}

	// Redis
	fmt.Println("Connecting to Redis")
	Redis, err = createRedisClient()

	if err != nil {
		sentry.CaptureException(err)
	}

	// Novu
	fmt.Println("Connecting to Novu Notifications (Admin)")
	Novu = createNovuClient()
}

// Discord
func createDiscordSession() (*discordgo.Session, error) {
	dg, err := discordgo.New("Bot " + Config.DiscordToken)
	if err != nil {
		return nil, err
	}

	return dg, nil
}

// Database
func createDatabase() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(Context, Config.DatabaseURL)
	if err != nil {
		return nil, errors.New("Uh oh! Unable to connect to Database: " + err.Error())
	}

	return pool, nil
}

// Redis
func createRedisClient() (*redis.Client, error) {
	opt, err := redis.ParseURL(Config.RedisURL)
	if err != nil {
		return nil, errors.New("Uh oh! Failed to connect to Redis: " + err.Error())
	}
	client := redis.NewClient(opt)

	return client, nil
}

// Novu
func createNovuClient() *novu.APIClient {
	return novu.NewAPIClient(Config.NovuAPIKey, &novu.Config{})
}
