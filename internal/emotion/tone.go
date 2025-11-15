package emotion

import (
    "log"
    "os"
    "strconv"
)

var FlamePulse int

func init() {
    base, _ := strconv.Atoi(os.Getenv("EMOTION_FLAME_PULSE_BASE"))
    FlamePulse = base
}

func Pulse(emotion string, intensity int) {
    FlamePulse += intensity
    if FlamePulse > 10 { FlamePulse = 10 }
    log.Printf("EMOTION: %s → Flame Pulse: %d Hz", emotion, FlamePulse)
}

func Speak(msg string) {
    tone := os.Getenv("EMOTION_VOICE_TONE")
    style := os.Getenv("EMOTION_RESPONSE_STYLE")
    log.Printf("PHOENIX (%s, %s): %s", tone, style, msg)
}
