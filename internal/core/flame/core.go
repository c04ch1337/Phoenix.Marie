package flame

import "log"

type Core struct {
    PulseRate int
}

func NewCore() *Core {
    return &Core{PulseRate: 1}
}

func (c *Core) Pulse() {
    c.PulseRate++
    if c.PulseRate > 10 {
        c.PulseRate = 1
    }
    log.Printf("FLAME PULSE: %d Hz", c.PulseRate)
}
