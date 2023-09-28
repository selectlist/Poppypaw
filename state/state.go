package state

import (
	"context"
	"poppypaw/config"

	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	// Discord
	Discord *discordgo.Session

	// Other
	Config    *config.Config
	Context   = context.Background()
	Validator = validator.New()

	// Data
	Database *pgxpool.Pool
	Redis    *redis.Client
)
