package routing

import (
    "encoding/json"
    "os"
    "fmt"
    "strings"
    "net/http"
    "html/template"
    "log"
    "bytes"
    "mime"

    "webservice/configuration"
    "webservice/state"

    f "github.com/gofiber/fiber/v2"
)


func SetRoutes( router *f.App, config *configuration.Config, store state.Store, healthiness *bool ) error {

    indexHtmlTemplate, err := template.New( "index" ).Parse( indexHtml )
    if err != nil {
        return err
    }

    if config.LogLevel == "debug" {
        router.All( "*", func( c *f.Ctx ) error {
            log.Printf( "%s %s  mime:%s  agent:%s",
                c.Method(),
                c.Path(),
                c.Get( f.HeaderContentType ),
                c.Get( f.HeaderUserAgent ),
            )
            return c.Next()
        })
    }


    router.Get( "/", func( c *f.Ctx ) error {
        headers := c.GetReqHeaders()
        acceptHeader := strings.Join( headers[ "Accept" ], " " )
        if ! strings.Contains( acceptHeader , "html" ) {
            c.Set( "Content-Type", "text/plain; charset=utf-8" )
            return c.SendString( "Hello, World!" )
        }

        data := indexHtmlData{
            Version: config.Version,
            Color: "",
        }

        buffer := &bytes.Buffer{}
        err := indexHtmlTemplate.Execute( buffer, data )
        if err != nil {
            return err
        }

        c.Set( "Content-Type", "text/html; charset=utf-8" )
        return c.Send( buffer.Bytes() )
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


    statePathGroup := router.Group( "/state" )


    statePathGroup.Get( "/:name", func( c *f.Ctx ) error {
        existingItem, err := store.Fetch( c.Params( "name" ) )
        if err != nil {
            c.Status( http.StatusInternalServerError )
            return c.Send( nil )
        }

        if existingItem == nil {
            return c.SendStatus( http.StatusNotFound )
        }

        c.Set( "Content-Type", existingItem.MimeType() )
        return c.Send( existingItem.Data() )
    })


    statePathGroup.Put( "/:name", func( c *f.Ctx ) error {
        contentType := strings.Clone( c.Get( "Content-Type" ) )
        _, _, err := mime.ParseMediaType( contentType )
        if err != nil {
            c.Status( http.StatusBadRequest )
            return c.SendString(
                fmt.Sprintf( "Invalid MIME type: %s", contentType ),
            )
        }

        name := strings.Clone( c.Params( "name" ) )
        existingItem, err := store.Fetch( name )
        if err != nil {
            c.Status( http.StatusInternalServerError )
            return c.Send( nil )
        }

        if existingItem != nil {
            if bytes.Equal( existingItem.Data(), c.Body() ) &&
               existingItem.MimeType() == contentType {
                c.Set( "Content-Type", "text/plain; charset=utf-8" )
                c.Status( http.StatusOK )
                return c.SendString( "Resource not changed" )
            }
            c.Status( http.StatusNoContent )
        } else {
            c.Status( http.StatusCreated )
        }

        newItem := state.NewItem(
            name,
            contentType,
            c.Body(),
        )

        if err = store.Add( newItem ); err != nil {
            c.Status( http.StatusInternalServerError )
            return c.Send( nil )
        }

        c.Set( "Content-Location", c.Path() )
        return c.Send( nil )
    })


    statePathGroup.Delete( "/:name", func( c *f.Ctx ) error {
        name := strings.Clone( c.Params( "name" ) )
        existingItem, err := store.Fetch( name )
        if err != nil {
            return c.SendStatus( http.StatusInternalServerError )
        }

        if existingItem == nil {
            return c.SendStatus( http.StatusNotFound )
        }

        if err = store.Remove( name ); err != nil {
            return c.SendStatus( http.StatusInternalServerError )
        }

        return c.SendStatus( http.StatusNoContent )
    })


    statePathGroup.Use( "*", func( c *f.Ctx ) error {
        if method := c.Method(); method == "OPTIONS" {
            c.Set( "Allow", "GET, PUT, DELETE, OPTIONS" )
            return c.SendStatus( http.StatusNoContent )
        }

        return c.SendStatus( http.StatusNotFound )
    })


    router.Get( "/states", func( c *f.Ctx ) error {
        states, err := store.Show()
        if err != nil {
            return c.SendStatus( http.StatusInternalServerError )
        }

        const pathPrefix string = "/state"
        paths := make ( []string, len( states ) )
        for i, state := range states {
            paths[ i ] = fmt.Sprintf( "%s/%s", pathPrefix, state )
        }

        headers := c.GetReqHeaders()
        acceptHeader := strings.Join( headers[ "Accept" ], " " )
        var response string
        if strings.Contains( acceptHeader, "json" ) {
            c.Set( "Content-Type", "application/json; charset=utf-8" )
            resJson, err := json.Marshal( paths )
            if err != nil {
                return err
            }
            response = string( resJson )
        } else {
            c.Set( "Content-Type", "text/plain; charset=utf-8" )
            response = strings.Join( paths, "\n" )
        }

        c.Status( http.StatusOK )
        return c.SendString( response )
    })


    router.Use( func( c *f.Ctx ) error {
        return c.SendStatus( http.StatusTeapot )
    })


    return nil
}
