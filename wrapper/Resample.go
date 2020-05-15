package wrapper

import (
	"fmt"
	"github.com/giorgisio/goav/avutil"
	"github.com/giorgisio/goav/swresample"
	"unsafe"
)

type Resample struct {
	swr              *swresample.Context
	frame            *avutil.Frame
	initialized      bool
	targetSampleRate int
	targetLayout     int
	targetSampleFmt  avutil.SampleFormat
}

func NewResample(rate int, channelLayout string, sampleFormat string) *Resample {

	layout := avutil.AvGetChannelLayout(channelLayout)
	format := avutil.SampleFormat(avutil.AvGetSampleFormat(sampleFormat))

	if layout == 0 || format == avutil.AV_SAMPLE_FMT_NONE {
		return nil
	}

	return &Resample{
		targetLayout:     layout,
		targetSampleRate: rate,
		targetSampleFmt:  format,
	}
}

func (r *Resample) reallocFrame(in *avutil.Frame) error {
	aframe := avutil.GetFrameAudioInfo(in)

	nchannels := avutil.AvGetNumberOfChannels(aframe.ChannelLayout)

	///
	if ret := avutil.AvAllocSamples(r.frame, nchannels, aframe.Samples, r.targetSampleFmt, 0); ret <= 0 {
		return fmt.Errorf("can not allocate samples %v", avutil.ErrorFromCode(ret))
	}

	aframe.ChannelLayout = r.targetLayout
	aframe.Format = r.targetSampleFmt
	aframe.SampleRate = r.targetSampleRate
	//aframe.Samples = aframe.Samples

	avutil.SetFrameAudioInfo(aframe, r.frame)

	return nil
}

func (r *Resample) init(in *avutil.Frame) error {

	// init
	r.swr = swresample.SwrAlloc()
	if r.swr == nil {
		return fmt.Errorf("Can not allocate SWR")
	}

	aframe := avutil.GetFrameAudioInfo(in)

	// setting up conversion
	r.swr.SwrSetOptionInt("in_channel_layout", aframe.ChannelLayout)
	r.swr.SwrSetOptionInt("in_sample_rate", aframe.SampleRate)
	r.swr.SwrSetSampleFmt("in_sample_fmt", aframe.Format)

	r.swr.SwrSetOptionInt("out_channel_layout", r.targetLayout)
	r.swr.SwrSetOptionInt("out_sample_rate", r.targetSampleRate)
	r.swr.SwrSetSampleFmt("out_sample_fmt", r.targetSampleFmt)

	r.frame = avutil.AvFrameAlloc()
	if r.frame == nil {
		panic(fmt.Errorf("can not allocate Audio Out"))
	}

	if ret := r.swr.SwrInit(); ret != 0 {
		return fmt.Errorf("can not init Audio Out: %v", avutil.ErrorFromCode(ret))
	}

	return r.reallocFrame(in)
}

func (r *Resample) Resample(in *avutil.Frame) (out *avutil.Frame, err error) {

	if !r.initialized {
		err := r.init(in)
		if err != nil {
			return nil, err
		}
		r.initialized = true
	}

	if ret := r.swr.SwrConvertFrame((*swresample.Frame)(unsafe.Pointer(r.frame)),
		(*swresample.Frame)(unsafe.Pointer(in))); ret != 0 {
		return nil, fmt.Errorf("can not convert frame %v", avutil.ErrorFromCode(ret))
	}

	return r.frame, nil
}
