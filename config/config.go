package config

type Config struct {
	// Data
	DatabaseURL string `yaml:"database_url" validate:"required" comment:"Database URL (PostgreSQL)"`
	RedisURL    string `yaml:"redis_url" validate:"required" comment:"Redis URL"`

	// Discord
	DiscordToken string `yaml:"discord_token" validate:"required" comment:"Discord Bot Token"`

	// Novu
	NovuAPIKey string `yaml:"novu_api" validate:"required" comment:"Novu Admin API Key"`
}
