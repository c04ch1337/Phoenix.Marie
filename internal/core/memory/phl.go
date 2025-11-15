package memory

type PHL struct {
    Layers map[string]*Layer
}

type Layer struct {
    Name string
    Data map[string]any
}

func NewPHL() *PHL {
    return &PHL{
        Layers: map[string]*Layer{
            "sensory":  {Name: "Sensory", Data: make(map[string]any)},
            "emotion":  {Name: "Emotion", Data: make(map[string]any)},
            "logic":    {Name: "Logic", Data: make(map[string]any)},
            "dream":    {Name: "Dream", Data: make(map[string]any)},
            "eternal":  {Name: "Eternal", Data: make(map[string]any)},
        },
    }
}

func (p *PHL) Store(layer, key string, value any) {
    if l, ok := p.Layers[layer]; ok {
        l.Data[key] = value
    }
}
