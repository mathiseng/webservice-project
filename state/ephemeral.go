package state

import (
    "errors"
    "sync"
)


type Ephemeral struct {
    store map[ string ] Item
    mux sync.Mutex
}


func NewEphemeralStore() *Ephemeral {
    return &Ephemeral{
        store: map[ string ] Item {},
        mux: sync.Mutex{},
    }
}


func ( e *Ephemeral ) Add( i Item ) error {
    if e.store == nil {
        return errors.New( "ephemeral storage not available" )
    }

    name := i.Name()

    e.mux.Lock()
    e.store[ name ] = i
    e.mux.Unlock()

    return nil
}


func ( e *Ephemeral ) Remove( name string ) error {
    if e.store == nil {
        return errors.New( "ephemeral storage not available" )
    }

    e.mux.Lock()
    delete( e.store, name )
    e.mux.Unlock()

    return nil
}


func ( e *Ephemeral ) Fetch( name string ) ( *Item, error ) {
    if e.store == nil {
        return nil, errors.New( "ephemeral storage not available" )
    }

    e.mux.Lock()
    item, found := e.store[ name ]
    e.mux.Unlock()

    if !found {
        return nil, nil
    }
    return &item, nil
}


func ( e *Ephemeral ) List() ( []string, error ) {
    if e.store == nil {
        return nil, errors.New( "ephemeral storage not available" )
    }

    e.mux.Lock()
    names := make( []string, 0, len( e.store ) )
    for _, item := range e.store {
        names = append( names, item.Name() )
    }
    e.mux.Unlock()

    return names, nil
}


func ( e *Ephemeral ) Disconnect() error {
    e.store = nil
    return nil
}
