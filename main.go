package main

import (
    "context"
    "fmt"
    "log"
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
        log.Fatalf( "HTTP server failed to start: %v", err )
    }

    server := fiber.New( fiber.Config{
        AppName: "webservice",
        DisableStartupMessage: config.Environment != "development",
    })

    store := state.NewEphemeralStore()

    var isHealthy = false

    routing.SetRoutes( server, config, store, &isHealthy )

    go func(){
        err := server.Listen( fmt.Sprintf( "%s:%d", config.Host, config.Port ) )
        if err != nil {
            log.Fatalf( "HTTP server failed to start: %v", err )
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
