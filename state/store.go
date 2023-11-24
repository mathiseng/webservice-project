package state



type Store interface {
    Add( i *Item ) error
    Remove( name string ) error
    Fetch( name string ) ( *Item, error )
    Show() ( []string, error )

    Disconnect() error
}
