package security

type ORCHDNA struct {
    ID     string
    Locked bool
}

func NewORCHDNA(id string) *ORCHDNA {
    return &ORCHDNA{ID: id, Locked: true}
}
