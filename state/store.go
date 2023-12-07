package state



type Store interface {
    Add( i Item ) error
    Remove( name string ) error
    Fetch( name string ) ( *Item, error )
    List() ( []string, error )

    Disconnect() error
}
