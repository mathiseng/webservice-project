package routing

import (
    "encoding/json"
    "os"
    "fmt"
    "net/http"

    "webservice/configuration"

    f "github.com/gofiber/fiber/v2"
)


func SetRoutes( router *f.App, config *configuration.Config, healthiness *bool ){

    router.Get( "/", func( c *f.Ctx ) error {
        return c.SendString( "Hello, World!" )
    })


    router.Get( "/health", func( c *f.Ctx ) error {
        type response struct {
            Status  string  `json:"status"  validate:"oneof=pass fail"`
        }

        c.Set( "Content-Type", "application/health+json; charset=utf-8" )

        var res *response
        if *healthiness == false {
            res = &response{
                Status: "fail",
            }
            c.Status( http.StatusServiceUnavailable )
        } else {
            res = &response{
                Status: "pass",
            }
            c.Status( http.StatusOK )
        }

        resJson, err := json.Marshal( res )
        if err != nil {
            return err
        }
        return c.SendString( string( resJson ) )
    })


    router.Get( "/env", func( c *f.Ctx ) error {
        c.Type( "txt", "utf-8" )

        if config.Environment == "production" {
            c.Status( http.StatusForbidden )
            return nil
        }

        for _, envVar := range os.Environ() {
            _, err := c.WriteString( fmt.Sprintln( envVar ) )
            if err != nil {
                c.Status( http.StatusInternalServerError )
                return err
            }
        }
        c.Status( http.StatusOK )

        return nil
    })


    router.Use( func( c *f.Ctx ) error {
        return c.SendStatus( http.StatusTeapot )
    })
}
