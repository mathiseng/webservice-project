package state

import (
    "testing"
    "mime"
    "sync"

    "github.com/stretchr/testify/assert"
)


var testItems = [] Item {
    NewItem( "foo", "bar", []byte( "fasel" ) ),
    NewItem( "qwertyASDFGH", mime.TypeByExtension( ".html" ), []byte{} ),
    Item{
        name: "Som!_🎵nam3",
        mimeType: "any kind of string",
        data: []byte{ 1, 2, 3, 4, 5, 6, 7, 8 },
    },
}


func TestEphemeralAdd( t *testing.T ){
    es := NewEphemeralStore()

    wg := &sync.WaitGroup{}
    for _, item := range testItems {
        wg.Add( 1 )
        go func( i Item ){
            defer wg.Done()
            es.Add( i )
        }( item )
    }
    wg.Wait()

    assert.Len( t, es.store, len( testItems ) )
}
