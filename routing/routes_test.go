package routing

import (
    "fmt"
    "io"
    "os"
    "time"
    "math/rand"
    "testing"
    "net/http"
    ht "net/http/httptest"

    f "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"

    "webservice/configuration"
)


func setup() ( *f.App, *configuration.Config ){
    os.Setenv( "ENV_NAME", "testing" )
    config, _ := configuration.New()

    server := f.New( f.Config{
        AppName: "test",
        DisableStartupMessage: false,
    })

    SetRoutes( server, config )

    return server, config
}


func bodyToString( body *io.ReadCloser ) ( string, error ) {
    defer ( *body ).Close()

    bodyBytes, err := io.ReadAll( *body )
    if err != nil {
        return "", err
    }
    return string( bodyBytes ), nil
}


func generateRandomNumberString() string {
    r := rand.New( rand.NewSource( time.Now().Unix() ) )
    randomNumber := r.Int63()
    return fmt.Sprintf( "%d", randomNumber )
}


func TestIndexRoute( t *testing.T ){
    router, _ := setup()

    req := ht.NewRequest( "GET", "/", nil )
    res, err := router.Test( req, -1 )

    bodyContent, err := bodyToString( &res.Body )
    assert.Nil( t, err )
    assert.Equal( t, "Hello, World!", bodyContent )
}


func TestEnvRoute( t *testing.T ){
    router, config := setup()

    envVarName := "TEST_ENV_VAR"
    envVarValue := generateRandomNumberString()

    os.Setenv( envVarName, envVarValue )

    req := ht.NewRequest( "GET", "/env", nil )
    res, err := router.Test( req, -1 )
    bodyContent, err := bodyToString( &res.Body )
    assert.Equal( t, http.StatusOK, res.StatusCode )
    assert.Nil( t, err )
    assert.Contains( t, bodyContent, fmt.Sprintf( "%s=%s", envVarName, envVarValue ) )

    ( *config ).Environment = "production"

    req = ht.NewRequest( "GET", "/env", nil )
    res, err = router.Test( req, -1 )
    assert.Equal( t, http.StatusForbidden, res.StatusCode )
}
