package state

import (
    "errors"
)


type Ephemeral struct {
    store map[ string ] *Item
}


func NewEphemeralStore() *Ephemeral {
    return &Ephemeral{
        store: map[ string ] *Item {},
    }
}


func ( e *Ephemeral ) Add( i *Item ) error {
    if e.store == nil {
        return errors.New( "ephemeral storage not available" )
    }

    name := i.Name()
    e.store[ name ] = i
    return nil
}


func ( e *Ephemeral ) Remove( name string ) error {
    if e.store == nil {
        return errors.New( "ephemeral storage not available" )
    }

    delete( e.store, name )
    return nil
}


func ( e *Ephemeral ) Fetch( name string ) ( *Item, error ) {
    if e.store == nil {
        return nil, errors.New( "ephemeral storage not available" )
    }

    item, found := e.store[ name ]
    if !found {
        return nil, nil
    }
    return item, nil
}


func ( e *Ephemeral ) Show() ( []string, error ) {
    if e.store == nil {
        return nil, errors.New( "ephemeral storage not available" )
    }

    names := make( []string, 0, len( e.store ) )
    for k := range e.store {
        names = append( names, k )
    }
    return names, nil
}


func ( e *Ephemeral ) Disconnect() error {
    e.store = nil
    return nil
}
