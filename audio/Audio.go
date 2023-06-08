package audio

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Sounds map[string]rl.Sound

func MainAudioLoop() {
	for {
		if !rl.IsSoundPlaying(Sounds["music_default"]) {
			rl.PlaySound(Sounds["music_default"])
			time.Sleep(time.Second * 27)
		}
	}
}

func InitAudio() {
	rl.InitAudioDevice()
	Sounds = make(map[string]rl.Sound)
	Sounds["music_default"] = rl.LoadSound("resources/audio/FranticLevel.wav")
}
