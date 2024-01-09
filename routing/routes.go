package routing

import (
    "encoding/json"
    "os"
    "fmt"
    "strings"
    "net/http"
    "html/template"
    log "log/slog"
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

    metricsTextTemplate, err := template.New( "metrics" ).Parse( metricsText )
    if err != nil {
        return err
    }

    if config.LogLevel == "debug" {
        router.All( "*", func( c *f.Ctx ) error {
            log.Debug(
                fmt.Sprintf( "%s %s  mime:%s  agent:%s",
                    c.Method(),
                    c.Path(),
                    c.Get( f.HeaderContentType, c.Get( f.HeaderAccept, "" ) ),
                    c.Get( f.HeaderUserAgent ),
                ),
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
            Color: config.FontColor,
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


    router.Get( "/metrics", func( c *f.Ctx ) error {
        headers := c.GetReqHeaders()
        acceptHeader := strings.Join( headers[ "Accept" ], " " )
        buffer := &bytes.Buffer{}

        if strings.Contains( acceptHeader , "json" ) {
            // FUTUREWORK: implement https://opentelemetry.io/docs/specs/otlp/#otlphttp
            return c.SendStatus( http.StatusNotAcceptable )
        } else {
            names, err := store.List()
            if err != nil {
                log.Debug( err.Error() )
                return c.SendStatus( http.StatusInternalServerError )
            }

            data := metricsTextData{
                Count: len( names ),
            }

            err = metricsTextTemplate.Execute( buffer, data )
            if err != nil {
                return err
            }

            c.Set( "Content-Type", "text/plain; charset=utf-8" )
            return c.Send( buffer.Bytes() )
        }
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
                log.Debug( err.Error() )
                c.Status( http.StatusInternalServerError )
                return err
            }
        }
        c.Status( http.StatusOK )

        return nil
    })


    statePathGroup := router.Group( "/state" )


    statePathGroup.Options( "/:name", func( c *f.Ctx ) error {
        name := strings.Clone( c.Params( "name" ) )
        existingItem, err := store.Fetch( name )
        if err != nil {
            log.Debug( err.Error() )
            return c.SendStatus( http.StatusInternalServerError )
        }

        if existingItem == nil {
            return c.SendStatus( http.StatusNotFound )
        }

        c.Set( "Allow", "OPTIONS, GET, PUT, DELETE, HEAD" )
        return c.SendStatus( http.StatusNoContent )
    })


    statePathGroup.Get( "/:name", func( c *f.Ctx ) error {
        existingItem, err := store.Fetch( c.Params( "name" ) )
        if err != nil {
            log.Debug( err.Error() )
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
            log.Debug( err.Error() )
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
            log.Debug( err.Error() )
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
            log.Debug( err.Error() )
            return c.SendStatus( http.StatusInternalServerError )
        }

        if existingItem == nil {
            return c.SendStatus( http.StatusNotFound )
        }

        if err = store.Remove( name ); err != nil {
            log.Debug( err.Error() )
            return c.SendStatus( http.StatusInternalServerError )
        }

        return c.SendStatus( http.StatusNoContent )
    })


    statePathGroup.Head( "/:name", func( c *f.Ctx ) error {
        name := strings.Clone( c.Params( "name" ) )
        existingItem, err := store.Fetch( name )
        if err != nil {
            log.Debug( err.Error() )
            return c.SendStatus( http.StatusInternalServerError )
        }

        if existingItem == nil {
            return c.SendStatus( http.StatusNotFound )
        }

        c.Set( "Content-Type", existingItem.MimeType() )
        c.Set( "Content-Length", fmt.Sprintf( "%d", len( existingItem.Data() ) ) )
        return c.SendStatus( http.StatusOK )
    })


    statePathGroup.Use( "*", func( c *f.Ctx ) error {
        return c.SendStatus( http.StatusNotFound )
    })


    router.Get( "/states", func( c *f.Ctx ) error {
        names, err := store.List()
        if err != nil {
            log.Debug( err.Error() )
            return c.SendStatus( http.StatusInternalServerError )
        }

        const pathPrefix string = "/state"
        paths := make( []string, len( names ) )
        for i, name := range names {
            paths[ i ] = fmt.Sprintf( "%s/%s", pathPrefix, name )
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
