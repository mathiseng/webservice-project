package routing

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
    "time"
    "math/rand"
    "encoding/json"
    "testing"
    "net/http"
    ht "net/http/httptest"

    f "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"

    "webservice/configuration"
    "webservice/state"
)


func setup() ( *f.App, *configuration.Config, state.Store, *bool ){
    os.Setenv( "ENV_NAME", "testing" )
    config, _ := configuration.New()

    server := f.New( f.Config{
        AppName: "test",
        DisableStartupMessage: false,
        BodyLimit: configuration.BODY_SIZE_LIMIT,
    })
    store := state.NewEphemeralStore()
    var isHealthy = true
    _ = SetRoutes( server, config, store, &isHealthy )

    return server, config, store, &isHealthy
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

func jsonToStringSlice( body *io.ReadCloser ) ( []string, error ) {
    defer ( *body ).Close()

    var data []string
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

func generateRandomBytes( length int ) []byte {
    r := rand.New( rand.NewSource( time.Now().Unix() ) )
    bytes := make( []byte, length )
    _, _ = r.Read( bytes )
    return bytes
}


func TestIndexRoute( t *testing.T ){
    router, _, _, _ := setup()

    req := ht.NewRequest( "GET", "/", nil )
    req.Header.Add( "Accept", "text/html" )
    res, _ := router.Test( req, -1 )
    bodyContent, err := bodyToString( &res.Body )
    assert.Nil( t, err )
    assert.Contains( t, bodyContent, "</html>" )
    assert.Contains( t, bodyContent, "<head>" )

    req = ht.NewRequest( "GET", "/", nil )
    res, _ = router.Test( req, -1 )
    bodyContent, err = bodyToString( &res.Body )
    assert.Nil( t, err )
    assert.Equal( t, "Hello, World!", bodyContent )
}


func TestHealthRoute( t *testing.T ){
    router, _, _, healthiness := setup()

    req := ht.NewRequest( "GET", "/health", nil )
    res, _ := router.Test( req, -1 )
    bodyContent, err := jsonToMap( &res.Body )
    status := bodyContent[ "status" ].( string )
    assert.Equal( t, http.StatusOK, res.StatusCode )
    assert.Nil( t, err )
    assert.Equal( t, "pass", status )

    *healthiness = false

    req = ht.NewRequest( "GET", "/health", nil )
    res, _ = router.Test( req, -1 )
    bodyContent, err = jsonToMap( &res.Body )
    status = bodyContent[ "status" ].( string )
    assert.Equal( t, http.StatusServiceUnavailable, res.StatusCode )
    assert.Nil( t, err )
    assert.Equal( t, "fail", status )
}


func TestEnvRoute( t *testing.T ){
    router, config, _, _ := setup()

    envVarName := "TEST_ENV_VAR"
    envVarValue := generateRandomNumberString()

    os.Setenv( envVarName, envVarValue )

    req := ht.NewRequest( "GET", "/env", nil )
    res, _ := router.Test( req, -1 )
    bodyContent, err := bodyToString( &res.Body )
    assert.Equal( t, http.StatusOK, res.StatusCode )
    assert.Nil( t, err )
    assert.Contains( t, bodyContent, fmt.Sprintf( "%s=%s", envVarName, envVarValue ) )

    ( *config ).Environment = "production"

    req = ht.NewRequest( "GET", "/env", nil )
    res, _ = router.Test( req, -1 )
    assert.Equal( t, http.StatusForbidden, res.StatusCode )
}


func TestState( t *testing.T ){
    router, _, store, _ := setup()

    const statePath1 = "/state/just-a-test"
    const statePath1Mime = "text/plain"
    const statePath1Body1 = "just a test body"
    const statePath1Body2 = "this body just changed"

    const statePath2 = "/state/another-test"
    const statePath2Mime = "application/octet-stream"
    const statePath2BodySize = 128
    statePath2Body := generateRandomBytes( statePath2BodySize )

    req := ht.NewRequest( "GET", statePath1, nil )
    res, _ := router.Test( req, -1 )
    assert.Equal( t, http.StatusNotFound, res.StatusCode )

    req = ht.NewRequest( "PUT", statePath1, nil )
    req.Header.Add( "Content-Type", "not a MIME type" )
    res, _ = router.Test( req, -1 )
    assert.Equal( t, http.StatusBadRequest, res.StatusCode )

    req = ht.NewRequest( "PUT", statePath1, strings.NewReader( statePath1Body1 ) )
    req.Header.Add( "Content-Type", statePath1Mime )
    res, _ = router.Test( req, -1 )
    assert.Equal( t, http.StatusCreated, res.StatusCode )

    req = ht.NewRequest( "PUT", statePath1, strings.NewReader( statePath1Body1 ) )
    req.Header.Add( "Content-Type", statePath1Mime )
    res, _ = router.Test( req, -1 )
    assert.Equal( t, http.StatusOK, res.StatusCode )

    req = ht.NewRequest( "PUT", statePath1, strings.NewReader( statePath1Body2 ) )
    req.Header.Add( "Content-Type", statePath1Mime )
    res, _ = router.Test( req, -1 )
    assert.Equal( t, http.StatusNoContent, res.StatusCode )

    req = ht.NewRequest( "GET", statePath1, nil )
    res, _ = router.Test( req, -1 )
    bodyContent, err := bodyToString( &res.Body )
    assert.Nil( t, err )
    assert.Equal( t, http.StatusOK, res.StatusCode )
    assert.Equal( t, statePath1Mime, res.Header[ "Content-Type" ][0] )
    assert.Equal( t, statePath1Body2, bodyContent )

    req = ht.NewRequest( "DELETE", statePath2, nil )
    res, _ = router.Test( req, -1 )
    assert.Equal( t, http.StatusNotFound, res.StatusCode )

    req = ht.NewRequest( "PUT", statePath2, bytes.NewReader( statePath2Body ) )
    req.Header.Add( "Content-Type", statePath2Mime )
    res, _ = router.Test( req, -1 )
    assert.Equal( t, http.StatusCreated, res.StatusCode )

    req = ht.NewRequest( "HEAD", statePath2, nil )
    res, _ = router.Test( req, -1 )
    contentLength, err := strconv.ParseInt( res.Header[ "Content-Length" ][0], 10, 64 )
    assert.Nil( t, err )
    assert.Equal( t, http.StatusOK, res.StatusCode )
    assert.Equal( t, statePath2Mime, res.Header[ "Content-Type" ][0] )
    assert.Equal( t, int64( statePath2BodySize ), contentLength )
    assert.IsType( t, res.Body, http.NoBody )

    req = ht.NewRequest( "GET", statePath2, nil )
    res, _ = router.Test( req, -1 )
    bodyBytes, err := io.ReadAll( res.Body )
    assert.Equal( t, http.StatusOK, res.StatusCode )
    assert.Equal( t, statePath2Mime, res.Header[ "Content-Type" ][0] )
    assert.Equal( t, statePath2Body, bodyBytes )

    req = ht.NewRequest( "GET", "/states", nil )
    req.Header.Add( "Accept", "application/json" )
    res, _ = router.Test( req, -1 )
    states, err := jsonToStringSlice( &res.Body )
    assert.Nil( t, err )
    assert.Len( t, states, 2 )
    assert.Contains( t, states, statePath1 )
    assert.Contains( t, states, statePath2 )

    req = ht.NewRequest( "DELETE", statePath1, nil )
    res, _ = router.Test( req, -1 )
    assert.Equal( t, http.StatusNoContent, res.StatusCode )

    req = ht.NewRequest( "GET", "/states", nil )
    res, _ = router.Test( req, -1 )
    statesPlain, err := bodyToString( &res.Body )
    assert.Nil( t, err )
    assert.NotContains( t, statesPlain, statePath1 )
    assert.Contains( t, statesPlain, statePath2 )

    err = store.Disconnect()
    req = ht.NewRequest( "GET", statePath2, nil )
    res, _ = router.Test( req, -1 )
    assert.Nil( t, err )
    assert.Equal( t, http.StatusInternalServerError, res.StatusCode )
}
