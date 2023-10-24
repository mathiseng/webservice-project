package routing

import (
    "os"
    "fmt"
    "net/http"

    "webservice/configuration"

    f "github.com/gofiber/fiber/v2"
)


func SetRoutes( router *f.App, config *configuration.Config ){

    router.Get( "/", func( c *f.Ctx ) error {
        return c.SendString( "Hello, World!" )
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
