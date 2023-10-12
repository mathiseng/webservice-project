package configuration

import (
    "errors"
    "fmt"

    configParser "github.com/caarlos0/env/v9"
)


type Config struct {
    Environment string `env:"ENV_NAME" envDefault:"development"`
    Host        string `env:"HOST" envDefault:"127.0.0.1"`
    Port        int16  `env:"PORT" envDefault:"3000"`
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

    return &cfg, nil
}
