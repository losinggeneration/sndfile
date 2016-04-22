package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/losinggeneration/openal"
	"github.com/losinggeneration/sndfile"
	"golang.org/x/mobile/exp/audio/al"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal(os.Args[0] + " <filename>")
	}

	filename := os.Args[1]

	if err := al.OpenDevice(); err != nil {
		log.Fatal(err)
	}

	s := al.GenSources(1)[0]
	defer al.DeleteSources(s)
	b := al.GenBuffers(1)[0]
	defer al.DeleteBuffers(b)

	snd, err := sndfile.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer snd.Close()

	buff, err := snd.ReadFrames(snd.Info.Frames)
	if err != nil {
		log.Fatal(err)
	}

	var f uint32
	if snd.Info.Channels == 2 {
		// 16 is chosen because ReadFrames defaults to reading int16's
		f = openal.FORMAT_STEREO16
	} else {
		f = openal.FORMAT_MONO16
	}

	b.BufferData(f, sndfile.Int16ToByte(buff), int32(snd.Info.SampleRate))

	s.QueueBuffers(b)
	s.Seti(openal.LOOPING, int32(0))
	s.SetGain(1.0)
	al.PlaySources(s)

	t := float64(snd.Info.Frames) / float64(snd.Info.SampleRate)
	hour, minute, second := time.Now().Add(time.Duration(t) * time.Second).Clock()
	fmt.Printf("Playing now. It will play for %.03f seconds and should be done playing at %02d:%02d:%02d\n", t, hour, minute, second)
	<-time.After(time.Duration(t) * time.Second)
}
