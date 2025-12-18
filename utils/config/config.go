package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/yorukot/knocker/models"
)

type AppEnv string

const (
	AppEnvDev  AppEnv = "dev"
	AppEnvProd AppEnv = "prod"
)

// EnvConfig holds all environment variables for the application
type EnvConfig struct {
	AppEnv       AppEnv   `env:"APP_ENV" envDefault:"prod"`
	AppName      string   `env:"APP_NAME" envDefault:"knocker"`
	AppMachineID int16    `env:"APP_MACHINE_ID" envDefault:"1"`
	AppPort      string   `env:"APP_PORT" envDefault:"8000"`
	AppRegion    string   `env:"APP_REGION" envDefault:"TW-Taipei"`
	AppRegions   []string `env:"APP_REGIONS" envDefault:"TW-Taipei" envSeparator:","`

	// Security Settings
	JWTSecretKey   string `env:"JWT_SECRET_KEY,required" envDefault:"change_me_to_a_secure_key"`
	FrontendDomain string `env:"FRONTEND_DOMAIN" envDefault:"localhost"`

	// PostgreSQL Settings
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSL_MODE,required"`

	// Redis/Dragonfly Settings
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`

	// Optional Settings
	OAuthStateExpiresAt   int `env:"OAUTH_STATE_EXPIRES_AT" envDefault:"600"`        // 10 minutes
	AccessTokenExpiresAt  int `env:"ACCESS_TOKEN_EXPIRES_AT" envDefault:"900"`       // 15 minutes
	RefreshTokenExpiresAt int `env:"REFRESH_TOKEN_EXPIRES_AT" envDefault:"31536000"` // 365 days

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID,required"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET,required"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL,required"`

	// Region settings (IDs are the primary source; names kept for legacy uses)
	AppRegionID  int64   `env:"APP_REGION_ID" envDefault:"1"`
	AppRegionIDs []int64 `env:"APP_REGION_IDS" envDefault:"1" envSeparator:","`
}

var (
	appConfig    *EnvConfig
	once         sync.Once
	regionsOnce  sync.Once
	regionsErr   error
	regionsByID  map[int64]models.Region
	regionsByKey map[string]models.Region
)

// loadConfig loads and validates all environment variables
func loadConfig() (*EnvConfig, error) {
	cfg := &EnvConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// InitConfig initializes the config only once
func InitConfig() (*EnvConfig, error) {
	var err error
	once.Do(func() {
		appConfig, err = loadConfig()
	})
	return appConfig, err
}

// Env returns the config. Panics if not initialized.
func Env() *EnvConfig {
	if appConfig == nil {
		panic("config not initialized — call InitConfig() first")
	}
	return appConfig
}

// InitRegionConfig loads regions from the database once and caches them.
// region_name matches legacy AppRegion/AppRegions, IDs match AppRegionID/AppRegionIDs.
func InitRegionConfig(pool *pgxpool.Pool) ([]models.Region, error) {
	ctx := context.Background()
	regionsOnce.Do(func() {
		rows, err := pool.Query(ctx, `SELECT id, region_name, region_display_name FROM regions ORDER BY id`)
		if err != nil {
			regionsErr = fmt.Errorf("list regions: %w", err)
			return
		}
		defer rows.Close()

		byID := make(map[int64]models.Region)
		byKey := make(map[string]models.Region)

		for rows.Next() {
			var r models.Region
			if err := rows.Scan(&r.ID, &r.Name, &r.DisplayName); err != nil {
				regionsErr = fmt.Errorf("scan region: %w", err)
				return
			}
			byID[r.ID] = r
			if key := r.Name; key != "" {
				byKey[key] = r
			}
		}

		if err := rows.Err(); err != nil {
			regionsErr = fmt.Errorf("iterate regions: %w", err)
			return
		}

		regionsByID = byID
		regionsByKey = byKey
	})

	return Regions(), regionsErr
}

// Regions returns cached regions sorted by ID; empty if not initialized.
func Regions() []models.Region {
	result := make([]models.Region, 0, len(regionsByID))
	for _, r := range regionsByID {
		result = append(result, r)
	}
	return result
}

// RegionByID looks up a cached region by ID.
func RegionByID(id int64) models.Region {
	r, _ := regionsByID[id]
	return r
}

// RegionByName looks up a cached region by its region_name (case-insensitive).
func RegionByName(name string) models.Region {
	r, _ := regionsByKey[name]
	return r
}
