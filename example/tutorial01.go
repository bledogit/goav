package main

// tutorial01.c
// Code based on a tutorial at http://dranger.com/ffmpeg/tutorial01.html

// A small sample program that shows how to use libavformat and libavcodec to
// read video from a file.
//
// Use
//
// gcc -o tutorial01 tutorial01.c -lavformat -lavcodec -lswscale -lz
//
// to build (assuming libavformat and libavcodec are correctly installed
// your system).
//
// Run using
//
// tutorial01 myvideofile.mpg
//
// to write the first five frames from "myvideofile.mpg" to disk in PPM
// format.
import (
	"fmt"
	"os"
	"unsafe"

	"gitlab.com/nielsen-media/eng/reference/commons/goav/avcodec"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avformat"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avutil"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a movie file")
		os.Exit(1)
	}

	// Open video file
	pFormatContext := avformat.AvformatAllocContext()
	if avformat.AvformatOpenInput(&pFormatContext, os.Args[1], nil, nil) != 0 {
		fmt.Printf("Unable to open file %s\n", os.Args[1])
		os.Exit(1)
	}

	// Retrieve stream information
	if pFormatContext.AvformatFindStreamInfo(nil) < 0 {
		fmt.Println("Couldn't find stream information")
		os.Exit(1)
	}

	// Dump information about file onto standard error
	pFormatContext.AvDumpFormat(0, os.Args[1], 0)

	// Find the first video stream
	for i := 0; i < int(pFormatContext.NbStreams()); i++ {
		theType := pFormatContext.Streams()[i].CodecParameters().AvCodecGetType()
		switch theType {
		case avformat.AVMEDIA_TYPE_AUDIO: //.AVMEDIA_TYPE_VIDEO:
			fmt.Println("Found Audio")
			// Get a pointer to the codec context for the video stream
			pCodecCtxOrig := pFormatContext.Streams()[i].Codec()
			// Find the decoder for the video stream
			pCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(pCodecCtxOrig.GetCodecId()))
			if pCodec == nil {
				fmt.Println("Unsupported codec!")
				os.Exit(1)
			}
			// Copy context
			pCodecCtx := pCodec.AvcodecAllocContext3()
			if pCodecCtx.AvCodecCopyContext((*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig))) != 0 {
				fmt.Println("Couldn't copy codec context")
				os.Exit(1)
			}

			// Open codec
			if pCodecCtx.AvcodecOpen2(pCodec, nil) < 0 {
				fmt.Println("Could not open codec")
				os.Exit(1)
			}

			// Allocate video frame
			pFrame := avutil.AvFrameAlloc()

			// Read frames and save first five frames to disk
			//frameNumber := 1
			packet := avcodec.AvPacketAlloc()
			for pFormatContext.AvReadFrame(packet) >= 0 {

				fmt.Println(packet)
				// Is this a packet from the video stream?
				if packet.StreamIndex() == i {
					// Decode video frame
					response := pCodecCtx.AvcodecSendPacket(packet)
					if response < 0 {
						fmt.Printf("Error while sending a packet to the decoder: %s\n", avutil.ErrorFromCode(response))
					}
					for response >= 0 {
						response = pCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(pFrame)))
						if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
							break
						} else if response < 0 {
							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
							return
						}

						fmt.Println(pFrame)

					}
				}

				// Free the packet that was allocated by av_read_frame
				packet.AvFreePacket()
			}

			// Free the YUV frame
			avutil.AvFrameFree(pFrame)

			// Close the codecs
			pCodecCtx.AvcodecClose()
			(*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig)).AvcodecClose()

			// Close the video file
			pFormatContext.AvformatCloseInput()

			// Stop after saving frames of first video straem
			break

		case avformat.AVMEDIA_TYPE_VIDEO:
			fmt.Println("Found Video")
		default:
			fmt.Println("Didn't find a video stream")
			os.Exit(1)
		}
	}
}
