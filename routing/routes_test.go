package routing

import (
    "fmt"
    "io"
    "os"
    "time"
    "math/rand"
    "encoding/json"
    "testing"
    "net/http"
    ht "net/http/httptest"

    f "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"

    "webservice/configuration"
)


func setup() ( *f.App, *configuration.Config, *bool ){
    os.Setenv( "ENV_NAME", "testing" )
    config, _ := configuration.New()

    server := f.New( f.Config{
        AppName: "test",
        DisableStartupMessage: false,
    })

    var isHealthy = true
    SetRoutes( server, config, &isHealthy )

    return server, config, &isHealthy
}


func bodyToString( body *io.ReadCloser ) ( string, error ) {
    defer ( *body ).Close()

    bodyBytes, err := io.ReadAll( *body )
    if err != nil {
        return "", err
    }
    return string( bodyBytes ), nil
}


func jsonToMap( body *io.ReadCloser ) ( map[string]interface{}, error ) {
    defer ( *body ).Close()

    var data map[string]interface{}
    bodyBytes, err := io.ReadAll( *body )
    if err != nil {
        return data, err
    }
    if err := json.Unmarshal( bodyBytes, &data ); err != nil {
       return data, err
    }
    return data, nil
}


func generateRandomNumberString() string {
    r := rand.New( rand.NewSource( time.Now().Unix() ) )
    randomNumber := r.Int63()
    return fmt.Sprintf( "%d", randomNumber )
}


func TestIndexRoute( t *testing.T ){
    router, _, _ := setup()

    req := ht.NewRequest( "GET", "/", nil )
    res, err := router.Test( req, -1 )

    bodyContent, err := bodyToString( &res.Body )
    assert.Nil( t, err )
    assert.Equal( t, "Hello, World!", bodyContent )
}


func TestHealthRoute( t *testing.T ){
    router, _, healthiness := setup()

    req := ht.NewRequest( "GET", "/health", nil )
    res, err := router.Test( req, -1 )
    bodyContent, err := jsonToMap( &res.Body )
    status := bodyContent[ "status" ].( string )
    assert.Equal( t, http.StatusOK, res.StatusCode )
    assert.Nil( t, err )
    assert.Equal( t, "pass", status )

    *healthiness = false

    req = ht.NewRequest( "GET", "/health", nil )
    res, err = router.Test( req, -1 )
    bodyContent, err = jsonToMap( &res.Body )
    status = bodyContent[ "status" ].( string )
    assert.Equal( t, http.StatusServiceUnavailable, res.StatusCode )
    assert.Nil( t, err )
    assert.Equal( t, "fail", status )
}


func TestEnvRoute( t *testing.T ){
    router, config, _ := setup()

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
