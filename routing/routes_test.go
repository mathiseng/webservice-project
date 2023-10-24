package routing

import (
    "io"
    "testing"
    ht "net/http/httptest"

    f "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
)


func setup() *f.App {
    router := f.New()

    SetRoutes( router )

    return router
}


func bodyToString( body *io.ReadCloser ) ( string, error ) {
    defer ( *body ).Close()

    bodyBytes, err := io.ReadAll( *body )
    if err != nil {
        return "", err
    }
    return string( bodyBytes ), nil
}


func TestIndexRoute( t *testing.T ){
    router := setup()

    req := ht.NewRequest( "GET", "/", nil )
    res, err := router.Test( req, -1 )

    bodyContent, err := bodyToString( &res.Body )
    assert.Nil( t, err )
    assert.Equal( t, "Hello World!", bodyContent )
}
