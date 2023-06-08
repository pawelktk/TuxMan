package audio

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Sounds map[string]rl.Sound

func MainAudio() {

	if !rl.IsSoundPlaying(Sounds["music_default"]) {
		rl.PlaySound(Sounds["music_default"])
	}

}

func InitAudio() {
	rl.InitAudioDevice()
	Sounds = make(map[string]rl.Sound)
	Sounds["music_default"] = rl.LoadSound("resources/audio/FranticLevel.wav")
}
