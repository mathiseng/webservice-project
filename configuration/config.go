package configuration

import (
    "errors"
    "fmt"

    configParser "github.com/caarlos0/env/v9"
)


const BODY_SIZE_LIMIT = 32 * 1024 * 1024    // 32 MB, in bytes


type Config struct {
    Version     string `env:"VERSION"   envDefault:"N/A"`

    LogLevel    string `env:"LOG_LEVE"  envDefault:"error"`

    Environment string `env:"ENV_NAME" envDefault:"development"`
    Host        string `env:"HOST" envDefault:"127.0.0.1"`
    Port        int16  `env:"PORT" envDefault:"3000"`

    DatabaseHost        string `env:"DB_HOST"       envDefault:""`
    DatabasePort        int16  `env:"DB_PORT"       envDefault:"6379"`
    DatabaseName        int    `env:"DB_NAME"       envDefault:"0"`
    DatabaseUsername    string `env:"DB_USERNAME"   envDefault:""`
    DatabasePassword    string `env:"DB_PASSWORD"   envDefault:""`
}


func New() ( *Config, error ){
    cfg := Config{}

    if err := configParser.Parse( &cfg ); err != nil {
        return nil, err
    }

    possibleEnvValues := map[ string ] bool {
        "development":  true,
        "testing":      true,
        "production":   true,
    }
    if _, ok := possibleEnvValues[ cfg.Environment ]; !ok {
        return nil, errors.New(
            fmt.Sprintf( "Invalid environment value: %s", cfg.Environment ),
        )
    }


    if cfg.Environment == "development" { cfg.LogLevel = "debug" }

    possibleLogLevels := map[ string ] bool {
        "error":    true,
        "debug":    true,
    }
    if _, ok := possibleLogLevels[ cfg.LogLevel ]; !ok {
        return nil, errors.New(
            fmt.Sprintf( "Invalid log level: %s", cfg.LogLevel ),
        )
    }


    return &cfg, nil
}
