package state


type Item struct {
    name        string
    mimeType    string
    data        []byte
}


func NewItem( name string, mimeType string, data []byte ) Item {
    return Item{
        name: name,
        mimeType: mimeType,
        data: data,
    }
}


func ( i *Item ) Name() string {
    return i.name
}

func ( i *Item ) MimeType() string {
    return i.mimeType
}

func ( i *Item ) Data() []byte {
    return i.data
}
