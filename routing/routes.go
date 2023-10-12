package routing

import (
    f "github.com/gofiber/fiber/v2"
)


func SetRoutes( router *f.App ){
    router.Get( "/", func( c *f.Ctx ) error {
        return c.SendString( "Hello World!" )
    })
}
