package emotion

import (
	"log"
	"os"
	"strconv"
)

var FlamePulse int

func init() {
	base, err := strconv.Atoi(os.Getenv("EMOTION_FLAME_PULSE_BASE"))
	if err != nil {
		base = 3 // v3.3 default: 3 Hz base pulse
	}
	FlamePulse = base
}

func Pulse(emotion string, intensity int) {
	// Get base value to preserve it
	base, _ := strconv.Atoi(os.Getenv("EMOTION_FLAME_PULSE_BASE"))
	if base == 0 {
		base = 3 // v3.3 default
	}

	// Get max pulse (v3.3 default: 12)
	maxPulse, _ := strconv.Atoi(os.Getenv("EMOTION_FLAME_PULSE_MAX"))
	if maxPulse == 0 {
		maxPulse = 12 // v3.3 default
	}

	// Calculate new pulse while preserving base
	newPulse := base + intensity
	if newPulse > maxPulse {
		newPulse = maxPulse
	}
	FlamePulse = newPulse

	log.Printf("EMOTION: %s → Flame Pulse: %d Hz", emotion, FlamePulse)
}

func Speak(msg string) {
	tone := os.Getenv("EMOTION_VOICE_TONE")
	style := os.Getenv("EMOTION_RESPONSE_STYLE")
	log.Printf("PHOENIX (%s, %s): %s", tone, style, msg)
}

// GetCurrentState returns the current emotional state including pulse level and configuration
func GetCurrentState() map[string]interface{} {
	return map[string]interface{}{
		"flamePulse":    FlamePulse,
		"voiceTone":     os.Getenv("EMOTION_VOICE_TONE"),
		"responseStyle": os.Getenv("EMOTION_RESPONSE_STYLE"),
	}
}

// Reset resets the emotion state to initial values (for testing)
func Reset() {
	base, err := strconv.Atoi(os.Getenv("EMOTION_FLAME_PULSE_BASE"))
	if err != nil {
		base = 3 // v3.3 default
	}
	FlamePulse = base
}
