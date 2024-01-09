package configuration

import (
    "errors"
    "fmt"
    "os"
    "log/slog"
    "unicode"
    fp "path/filepath"

    configParser "github.com/caarlos0/env/v9"
)


const BODY_SIZE_LIMIT = 32 * 1024 * 1024    // 32 MB, in bytes

var version string = "n/a"


type Config struct {
    Version     string

    FontColor   string `env:"FONT_COLOR"  envDefault:""`

    LogLevel    string `env:"LOG_LEVE"  envDefault:"error"`

    Environment     string `env:"ENV_NAME"  envDefault:"development"`
    Host            string `env:"HOST"      envDefault:"127.0.0.1"`
    Port            int16  `env:"PORT"      envDefault:"3000"`

    DatabaseHost        string `env:"DB_HOST"       envDefault:""`
    DatabasePort        int16  `env:"DB_PORT"       envDefault:"6379"`
    DatabaseName        int    `env:"DB_NAME"       envDefault:"0"`
    DatabaseUsername    string `env:"DB_USERNAME"   envDefault:""`
    DatabasePassword    string `env:"DB_PASSWORD"   envDefault:""`
}


func New() ( *Config, error ){
    cfg := Config{
        Version: version,
    }

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

    if _, err := cfg.GetLogLevel(); err != nil {
        return nil, err
    }

    if len( cfg.DatabaseHost ) >= 1 && len( cfg.DatabasePassword ) >= 2 {
        if ! fp.IsLocal( cfg.DatabasePassword ) && ! fp.IsAbs( cfg.DatabasePassword ) {
            return nil, errors.New(
                fmt.Sprintln( "Database password must be a file path" ),
            )
        }
        _, err := os.Stat( cfg.DatabasePassword )
        if err != nil {
            if errors.Is( err, os.ErrNotExist ){
                return nil, errors.New(
                    fmt.Sprintln( "Database password file does not exist" ),
                )
            }
            return nil, errors.New(
                fmt.Sprintln( "Database password file not accessible" ),
            )
        }
    }

    if len( cfg.FontColor ) >= 1 {
        if len( cfg.FontColor ) >= 21 {
            return nil, errors.New(
                fmt.Sprintln( "Font color too long" ),
            )
        }

        for _, r := range cfg.FontColor {
            if ! unicode.IsLetter( r ) {
                return nil, errors.New(
                    fmt.Sprintln( "Invalid character in font color" ),
                )
            }
        }
    }


    return &cfg, nil
}


func ( cfg *Config ) GetLogLevel() ( slog.Level, error ){
    possibleLogLevels := map[ string ] slog.Level {
        "error":    slog.LevelError,
        "debug":    slog.LevelDebug,
    }
    level, ok := possibleLogLevels[ cfg.LogLevel ]
    if !ok {
        return slog.LevelError, errors.New(
            fmt.Sprintf( "Invalid log level: %s", cfg.LogLevel ),
        )
    }else{
        return level, nil
    }
}