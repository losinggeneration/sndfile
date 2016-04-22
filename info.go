package sndfile

// #include <sndfile.h>
import "C"

// Info contains information about the file (in ReadMode) or how the data should be written
type Info struct {
	Frames     uint // Frames specifies how many frames a file has
	SampleRate int  // SampleRate is the file's sample rate
	Channels   int  // Channels is how many channels the file includes
	Format     int  // Format is the file's format
	Sections   int  // Sections is how many sections a file has
	Seekable   bool // Seekable is if the file is seekable or not

	info C.SF_INFO
}

func (i *Info) fillInfo() {
	i.Frames = uint(i.info.frames)
	i.SampleRate = int(i.info.samplerate)
	i.Channels = int(i.info.channels)
	i.Format = int(i.info.format)
	i.Sections = int(i.info.sections)
	i.Seekable = i.info.seekable > 0
}

func (i *Info) fillCInfo() {
	i.info.frames = C.sf_count_t(i.Frames)
	i.info.samplerate = C.int(i.SampleRate)
	i.info.channels = C.int(i.Channels)
	i.info.format = C.int(i.Format)
	i.info.sections = C.int(i.Sections)
	if i.Seekable {
		i.info.seekable = 1
	} else {
		i.info.seekable = 0
	}
}
