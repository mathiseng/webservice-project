package main

import (
    "context"
    "fmt"
    "log"
    "log/slog"
    "os"
    "os/signal"
    "syscall"
    "time"

    "webservice/configuration"
    "webservice/routing"
    "webservice/state"

    "github.com/gofiber/fiber/v2"
)


func main() {
    config, err := configuration.New()
    if err != nil {
        slog.Error( fmt.Sprintf( "HTTP server failed to start: %v", err ) )
        os.Exit( 1 )
    }

    level, _ := config.GetLogLevel()
    slog.SetDefault(
        slog.New(
            slog.NewTextHandler(
                os.Stdout,
                &slog.HandlerOptions{
                    Level: level,
                },
            ),
        ),
    )

    server := fiber.New( fiber.Config{
        AppName: "webservice",
        DisableStartupMessage: config.Environment != "development",
        BodyLimit: configuration.BODY_SIZE_LIMIT,
    })

    var store state.Store
    if len( config.DatabaseHost ) <= 0 {
        store = state.NewEphemeralStore()
    } else {
        store = state.NewPersistentStore( config )
    }

    var isHealthy = false

    err = routing.SetRoutes( server, config, store, &isHealthy )
    if err != nil {
        slog.Error( fmt.Sprintf( "HTTP server failed to start: %v", err ) )
        os.Exit( 1 )
    }

    go func(){
        err := server.Listen( fmt.Sprintf( "%s:%d", config.Host, config.Port ) )
        if err != nil {
            slog.Error( fmt.Sprintf( "HTTP server failed to start: %v", err ) )
            os.Exit( 1 )
        }
    }()

    osSignaling := make( chan os.Signal, 1 )
    signal.Notify( osSignaling, syscall.SIGHUP  )
    signal.Notify( osSignaling, syscall.SIGINT  )
    signal.Notify( osSignaling, syscall.SIGTERM )
    signal.Notify( osSignaling, syscall.SIGQUIT )

    shuttingDown := context.TODO()
    isHealthy = true
    if config.Environment != "development" {
        log.Println( "HTTP server started successfully" )
    }

    for {
        select {
        case <-osSignaling:
            isHealthy = false
            log.Println( "Gracefully shutting down HTTP server" )

            var concludeShutdown context.CancelFunc
            shuttingDown, concludeShutdown = context.WithTimeout(
                context.Background(),
                time.Second * 15,
            )
            err := server.ShutdownWithContext( shuttingDown )
            if err != nil {
                log.Printf( "HTTP server failed to shut down: %v", err )
            }
            err = store.Disconnect()
            if err != nil {
                log.Printf( "Store failed to disconnect: %v", err )
            }
            concludeShutdown()

        case <-shuttingDown.Done():
            return
        }
    }
}
