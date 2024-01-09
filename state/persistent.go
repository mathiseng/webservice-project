package state

import (
    "fmt"
    "runtime"
    "context"
    "time"
    "os"
    log "log/slog"

    "webservice/configuration"

    db "github.com/redis/go-redis/v9"
)



type Persistent struct {
    client      *db.Client
    ctx         context.Context
    timeout     time.Duration
}


func NewPersistentStore( c *configuration.Config ) *Persistent {
    content, err := os.ReadFile( c.DatabasePassword )
    if err != nil {
        log.Error( fmt.Sprintf( "Database password not able to be read: %v", err ) )
        os.Exit( 1 )
    }
    dbPassword := string( content )

    return &Persistent{
        client: db.NewClient( &db.Options{
            Addr: fmt.Sprintf( "%s:%d", c.DatabaseHost, c.DatabasePort ),
            Username: c.DatabaseUsername,
            Password: dbPassword,
            DB: c.DatabaseName,

            DialTimeout: time.Second * 3,
            ContextTimeoutEnabled: true,

            MaxRetries: 3,
            MinRetryBackoff: time.Second * 1,
            MaxRetryBackoff: time.Second * 2,

            PoolSize: 10 * runtime.NumCPU(),
            MaxActiveConns: 10 * runtime.NumCPU(),
        }),

        timeout: time.Second * 20,
    }
}


func ( e *Persistent ) Add( i Item ) error {
    ctx, cancel := context.WithTimeout( context.TODO(), e.timeout )
    defer cancel()

    name := i.Name()
    if err := e.client.HSet(
        ctx, name,
        "mime", i.MimeType(),
        "data", i.Data(),
    ).Err(); err != nil {
        return err
    }
    return nil
}


func ( e *Persistent ) Remove( name string ) error {
    ctx, cancel := context.WithTimeout( context.TODO(), e.timeout )
    defer cancel()

    if err := e.client.Del( ctx, name ).Err(); err != nil {
        return err
    }
    return nil
}


func ( e *Persistent ) Fetch( name string ) ( *Item, error ) {
    ctx, cancel := context.WithTimeout( context.TODO(), e.timeout )
    defer cancel()

    value, err := e.client.HGetAll( ctx, name ).Result()
    if err != nil {
        return nil, err
    }

    var item *Item = nil
    if len( value ) >= 1 {
        i := NewItem( name, value[ "mime" ], []byte( value[ "data" ] ) )
        item = &i
    }
    return item, nil
}


func ( e *Persistent ) List() ( []string, error ) {
    ctx, cancel := context.WithTimeout( context.TODO(), e.timeout )
    defer cancel()

    var names []string
    i := e.client.Scan( ctx, 0, "", 0 ).Iterator()
    for i.Next( ctx ){
        names = append( names, i.Val() )
    }
    if err := i.Err(); err != nil {
        return nil, err
    }
    return names, nil
}


func ( e *Persistent ) Disconnect() error {
    return e.client.Close()
}