package main

import "C"
import (
	"fmt"
	"github.com/giorgisio/goav/avutil"
	"github.com/giorgisio/goav/wrapper"
	"os"
	"unsafe"
)

func SaveAudioFrame(file *os.File, frame *avutil.Frame) {
	// Open file
	data := avutil.Data(frame)
	linesize := avutil.Linesize(frame)
	af := avutil.AvGetFrameAudioInfo(frame)
	channels := avutil.AvGetNumberOfChannels(af.ChannelLayout)
	buf := C.GoBytes(unsafe.Pointer(data[0]), C.int(af.Samples*channels*2))

	file.Write(buf)
	fmt.Println(linesize)
	// Write  data
}

func SavePictureFrame(frame *avutil.Frame, width, height, frameNumber int) {
	// Open file
	fileName := fmt.Sprintf("frame%d.ppm", frameNumber)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error Reading")
	}
	defer file.Close()

	// Write header
	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
	file.Write([]byte(header))

	line := avutil.Linesize(frame)
	data := avutil.Data(frame)

	// Write pixel data
	for y := 0; y < height; y++ {
		startPos := uintptr(unsafe.Pointer(data[0])) + uintptr(y)*uintptr(line[0])
		buf := C.GoBytes(unsafe.Pointer(startPos), C.int(width))
		file.Write(buf)
		file.Write(buf)
		file.Write(buf)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide a movie file")
		os.Exit(1)
	}

	demuxer := wrapper.Demuxer{}

	err := demuxer.InitWithFile(os.Args[1])
	check(err)

	decoder := wrapper.Decoder{}

	stream := 0
	err = decoder.Init(demuxer, stream)
	check(err)

	resample := wrapper.NewResample(44100, "5.1", avutil.AV_SAMPLE_FMT_S16)
	if resample == nil {
		panic("can not initialize resample")
	}

	pcmout, _ := os.Create("out.pcm")
	for {
		packet := demuxer.Demux()

		if packet != nil && packet.StreamIndex() == stream {
			frame, err := decoder.Decode(packet)
			check(err)

			if frame != nil {

				af := avutil.AvGetFrameAudioInfo(frame)
				fmt.Printf("frame %d, %+v\n", packet.StreamIndex(), af)

				frameout, err := resample.Resample(frame)
				check(err)

				SaveAudioFrame(pcmout, frameout)
			}
		}
	}
}
