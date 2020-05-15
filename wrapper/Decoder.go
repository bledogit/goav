package wrapper

import (
	"fmt"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avutil"
	"unsafe"
)

type Decoder struct {
	codec *avcodec.Context
	frame *avutil.Frame
}

func (d *Decoder) Init(demux Demuxer, stream int) error {

	if stream >= int(demux.context.NbStreams()) {
		return fmt.Errorf("stream does not exists %d", stream)
	}

	codecCtxOrig := demux.GetContext().Streams()[stream].Codec()

	d.frame = avutil.AvFrameAlloc()

	decoder := avcodec.AvcodecFindDecoder(avcodec.CodecId(codecCtxOrig.GetCodecId()))
	if decoder == nil {
		return fmt.Errorf("unsupported codec")
	}

	if d.codec = decoder.AvcodecAllocContext3(); d.codec == nil {
		return fmt.Errorf("can not allocate codec")
	}

	// Open codec
	if d.codec.AvcodecOpen2(decoder, nil) < 0 {
		return fmt.Errorf("Could not open codec")
	}

	return nil
}

func (d *Decoder) Decode(packet *avcodec.Packet) (*avutil.Frame, error) {
	response := d.codec.AvcodecSendPacket(packet)

	if response < 0 {
		return nil, fmt.Errorf("Error while sending a packet to the decoder: %v\n", avutil.ErrorFromCode(response))
	}
	for response >= 0 {

		response = d.codec.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(d.frame)))

		if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
			fmt.Printf("break\n")
			break
		} else if response < 0 {
			return nil, fmt.Errorf("Error while receiving a frame from the decoder: %v\n", avutil.ErrorFromCode(response))
		}

		return d.frame, nil
	}
	return nil, nil
}

func (d *Decoder) Close() {

}
