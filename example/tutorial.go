package main

import "C"
import (
	"fmt"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avcodec"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avutil"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/wrapper"
	"os"
	"unsafe"
)

func SaveAudioFrame(file *os.File, frame *avutil.Frame) {
	// Open file
	data := avutil.Data(frame)

	af := avutil.GetFrameAudioInfo(frame)
	channels := avutil.AvGetNumberOfChannels(af.ChannelLayout)
	buf := C.GoBytes(unsafe.Pointer(data[0]), C.int(af.Samples*int64(channels)*2))

	file.Write(buf)
	// linesize := avutil.Linesize(frame)
	// fmt.Println(linesize)
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

	demuxer := wrapper.NewDemuxer()
	err := demuxer.WithFile(os.Args[1])
	check(err)

	stream := 1
	decoder, err := wrapper.NewDecoder(demuxer, stream)
	check(err)

	resample := wrapper.NewResample(44100, "5.1", "s16")
	if resample == nil {
		panic("can not initialize resample")
	}

	pcmout, _ := os.Create("out.pcm")
	packet := avcodec.AvPacketAlloc()
	for {
		if demuxer.Demux(packet) < 0 {
			break
		}

		if packet != nil && packet.StreamIndex() == stream {
			frames, err := decoder.Decode(packet)
			check(err)

			for _, frame := range frames {

				af := avutil.GetFrameAudioInfo(frame)
				fmt.Printf("stream %v frame %+v\n", packet.StreamIndex(), af)

				frameout, err := resample.Resample(frame)
				check(err)

				SaveAudioFrame(pcmout, frameout)

				avutil.AvFrameFree(frame)

			}
		}
	}
	packet.AvPacketUnref()
	decoder.Close()
	demuxer.Close()
}
