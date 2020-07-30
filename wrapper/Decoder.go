package wrapper

import (
	"fmt"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avcodec"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avutil"
	"unsafe"
)

type Decoder struct {
	codec *avcodec.Context
}

func NewDecoder(demux *Demuxer, stream int) (*Decoder, error) {
	d := &Decoder{}

	if stream >= int(demux.context.NbStreams()) {
		return nil, fmt.Errorf("stream does not exists %d", stream)
	}

	codecCtxOrig := demux.GetContext().Streams()[stream].Codec()

	decoder := avcodec.AvcodecFindDecoder(avcodec.CodecId(codecCtxOrig.GetCodecId()))
	if decoder == nil {
		return nil, fmt.Errorf("unsupported codec")
	}

	d.codec = decoder.AvcodecAllocContext3()
	if d.codec == nil {
		return nil, fmt.Errorf("can not allocate codec")
	}

	pars := avcodec.AvCodecAllocParameters()

	if (*avcodec.Context)(unsafe.Pointer(codecCtxOrig)).AvcodecParametersFromContext(pars) != 0 {
		return nil, fmt.Errorf("error copying parameteres")
	}

	if d.codec.AvcodecParametersToContext(pars) != 0 {
		return nil, fmt.Errorf("error copying parameteres")
	}

	pars.AvcodecParametersFree()

	// Open codec
	if d.codec.AvcodecOpen2(decoder, nil) < 0 {
		return nil, fmt.Errorf("Could not open codec")
	}

	return d, nil
}

func (d *Decoder) Decode(packet *avcodec.Packet) ([]*avutil.Frame, error) {
	response := d.codec.AvcodecSendPacket(packet)
	frameOut := []*avutil.Frame{}

	if response < 0 {
		return nil, fmt.Errorf("Error while sending a packet to the decoder: %v\n", avutil.ErrorFromCode(response))
	}
	for response >= 0 {

		frame := avutil.AvFrameAlloc()
		response = d.codec.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(frame)))

		if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEINVAL || response == avutil.AvErrorEOF {
			break
		} else if response < 0 {
			return nil, fmt.Errorf("Error while receiving a frame from the decoder: %v\n", avutil.ErrorFromCode(response))
		}
		frameOut = append(frameOut, frame)

	}
	return frameOut, nil
}

func (d *Decoder) Close() {
	d.codec.AvcodecClose()
}
