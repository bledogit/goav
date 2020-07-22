// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// Giorgis (habtom@giorgis.io)

package avutil

/*
	#cgo pkg-config: libavutil
	#include <libavutil/frame.h>
    #include <libavutil/channel_layout.h>
	#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"image"
	"log"
	"unsafe"
)

type (
	AvBuffer            C.struct_AVBuffer
	AvBufferRef         C.struct_AVBufferRef
	AvBufferPool        C.struct_AVBufferPool
	Frame               C.struct_AVFrame
	AvFrameSideData     C.struct_AVFrameSideData
	AvFrameSideDataType C.enum_AVFrameSideDataType
	SampleFormat        C.enum_AVSampleFormat
)

const (
	AV_SAMPLE_FMT_NONE = -1
	AV_SAMPLE_FMT_U8   = 0 ///< unsigned 8 bits
	AV_SAMPLE_FMT_S16  = 1 ///< signed 16 bits
	AV_SAMPLE_FMT_S32  = 2 ///< signed 32 bits
	AV_SAMPLE_FMT_FLT  = 3 ///< float
	AV_SAMPLE_FMT_DBL  = 4 ///< double

	AV_SAMPLE_FMT_U8P  = 5  ///< unsigned 8 bits, planar
	AV_SAMPLE_FMT_S16P = 6  ///< signed 16 bits, planar
	AV_SAMPLE_FMT_S32P = 7  ///< signed 32 bits, planar
	AV_SAMPLE_FMT_FLTP = 8  ///< float, planar
	AV_SAMPLE_FMT_DBLP = 9  ///< double, planar
	AV_SAMPLE_FMT_S64  = 10 ///< signed 64 bits
	AV_SAMPLE_FMT_S64P = 11 ///< signed 64 bits, planar

	AV_SAMPLE_FMT_NB = 12 ///< Number of sample formats. DO NOT USE if linking dynamically
)

// AudioFrame exports AvFrame fields
type AudioFrame struct {
	Samples       int64
	SampleRate    int64
	Pts           int64
	ChannelLayout int64
	Format        SampleFormat
}

func AvprivFrameGetMetadatap(f *Frame) *Dictionary {
	return (*Dictionary)(unsafe.Pointer(f.metadata))
}

//Allocate an Frame and set its fields to default values.
func AvFrameAlloc() *Frame {
	return (*Frame)(unsafe.Pointer(C.av_frame_alloc()))
}

//Free the frame and any dynamically allocated objects in it, e.g.
func AvFrameFree(f *Frame) {
	C.av_frame_free((**C.struct_AVFrame)(unsafe.Pointer(&f)))
}

//Allocate new buffer(s) for audio or video data.
func AvFrameGetBuffer(f *Frame, a int) int {
	return int(C.av_frame_get_buffer((*C.struct_AVFrame)(unsafe.Pointer(f)), C.int(a)))
}

//Setup a new reference to the data described by an given frame.
func AvFrameRef(d, s *Frame) int {
	return int(C.av_frame_ref((*C.struct_AVFrame)(unsafe.Pointer(d)), (*C.struct_AVFrame)(unsafe.Pointer(s))))
}

//Create a new frame that references the same data as src.
func AvFrameClone(f *Frame) *Frame {
	return (*Frame)(C.av_frame_clone((*C.struct_AVFrame)(unsafe.Pointer(f))))
}

//Unreference all the buffers referenced by frame and reset the frame fields.
func AvFrameUnref(f *Frame) {
	cf := (*C.struct_AVFrame)(unsafe.Pointer(f))
	C.av_frame_unref(cf)
}

//Move everythnig contained in src to dst and reset src.
func AvFrameMoveRef(d, s *Frame) {
	C.av_frame_move_ref((*C.struct_AVFrame)(unsafe.Pointer(d)), (*C.struct_AVFrame)(unsafe.Pointer(s)))
}

//Check if the frame data is writable.
func AvFrameIsWritable(f *Frame) int {
	return int(C.av_frame_is_writable((*C.struct_AVFrame)(unsafe.Pointer(f))))
}

//Ensure that the frame data is writable, avoiding data copy if possible.
func AvFrameMakeWritable(f *Frame) int {
	return int(C.av_frame_make_writable((*C.struct_AVFrame)(unsafe.Pointer(f))))
}

//Copy only "metadata" fields from src to dst.
func AvFrameCopyProps(d, s *Frame) int {
	return int(C.av_frame_copy_props((*C.struct_AVFrame)(unsafe.Pointer(d)), (*C.struct_AVFrame)(unsafe.Pointer(s))))
}

//Get the buffer reference a given data plane is stored in.
func AvFrameGetPlaneBuffer(f *Frame, p int) *AvBufferRef {
	return (*AvBufferRef)(C.av_frame_get_plane_buffer((*C.struct_AVFrame)(unsafe.Pointer(f)), C.int(p)))
}

//Add a new side data to a frame.
func AvFrameNewSideData(f *Frame, d AvFrameSideDataType, s int) *AvFrameSideData {
	return (*AvFrameSideData)(C.av_frame_new_side_data((*C.struct_AVFrame)(unsafe.Pointer(f)), (C.enum_AVFrameSideDataType)(d), C.int(s)))
}

func AvFrameGetSideData(f *Frame, t AvFrameSideDataType) *AvFrameSideData {
	return (*AvFrameSideData)(C.av_frame_get_side_data((*C.struct_AVFrame)(unsafe.Pointer(f)), (C.enum_AVFrameSideDataType)(t)))
}

func Data(f *Frame) (data [8]*uint8) {
	for i := range data {
		data[i] = (*uint8)(f.data[i])
	}
	return
}

func Linesize(f *Frame) (linesize [8]int32) {
	for i := range linesize {
		linesize[i] = int32(f.linesize[i])
	}
	return
}

func Pts(f *Frame) int64 {
	return int64(f.pts)
}

//GetPicture creates a YCbCr image from the frame
func GetPicture(f *Frame) (img *image.YCbCr, err error) {
	// For 4:4:4, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/1.
	// For 4:2:2, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/2.
	// For 4:2:0, CStride == YStride/2 && len(Cb) == len(Cr) == len(Y)/4.
	// For 4:4:0, CStride == YStride/1 && len(Cb) == len(Cr) == len(Y)/2.
	// For 4:1:1, CStride == YStride/4 && len(Cb) == len(Cr) == len(Y)/4.
	// For 4:1:0, CStride == YStride/4 && len(Cb) == len(Cr) == len(Y)/8.

	w := int(f.linesize[0])
	h := int(f.height)
	r := image.Rectangle{image.Point{0, 0}, image.Point{w, h}}
	// TODO: Use the sub sample ratio from the input image 'f.format'
	img = image.NewYCbCr(r, image.YCbCrSubsampleRatio420)
	// convert the frame data data to a Go byte array
	img.Y = C.GoBytes(unsafe.Pointer(f.data[0]), C.int(w*h))

	wCb := int(f.linesize[1])
	if unsafe.Pointer(f.data[1]) != nil {
		img.Cb = C.GoBytes(unsafe.Pointer(f.data[1]), C.int(wCb*h/2))
	}

	wCr := int(f.linesize[2])
	if unsafe.Pointer(f.data[2]) != nil {
		img.Cr = C.GoBytes(unsafe.Pointer(f.data[2]), C.int(wCr*h/2))
	}
	return
}

// SetPicture sets the image pointer of |f| to the image pointers of |img|
func SetPicture(f *Frame, img *image.YCbCr) {
	d := Data(f)
	// l := Linesize(f)
	// FIXME: Save the original pointers somewhere, this is a memory leak
	d[0] = (*uint8)(unsafe.Pointer(&img.Y[0]))
	// d[1] = (*uint8)(unsafe.Pointer(&img.Cb[0]))
}

func GetPictureRGB(f *Frame) (img *image.RGBA, err error) {
	w := int(f.linesize[0])
	h := int(f.height)
	r := image.Rectangle{image.Point{0, 0}, image.Point{w, h}}
	// TODO: Use the sub sample ratio from the input image 'f.format'
	img = image.NewRGBA(r)
	// convert the frame data data to a Go byte array
	img.Pix = C.GoBytes(unsafe.Pointer(f.data[0]), C.int(w*h))
	img.Stride = w
	log.Println("w", w, "h", h)
	return
}

func AvSetFrame(f *Frame, w int, h int, pixFmt int) (err error) {
	f.width = C.int(w)
	f.height = C.int(h)
	f.format = C.int(pixFmt)
	if ret := C.av_frame_get_buffer((*C.struct_AVFrame)(unsafe.Pointer(f)), 32 /*alignment*/); ret < 0 {
		err = fmt.Errorf("Error allocating avframe buffer. Err: %v", ret)
		return
	}
	return
}

func AvFrameVideoInfo(f *Frame) (width int, height int, linesize [8]int32, data [8]*uint8) {
	width = int(f.linesize[0])
	height = int(f.height)
	for i := range linesize {
		linesize[i] = int32(f.linesize[i])
	}
	for i := range data {
		data[i] = (*uint8)(f.data[i])
	}
	// log.Println("Linesize is ", f.linesize, "Data is", data)
	return
}

// AvGetNumberOfChannels calls av_get_channel_layout_nb_channels
func AvGetNumberOfChannels(layout int64) int {
	return int(C.av_get_channel_layout_nb_channels(C.uint64_t(layout)))
}

func GetBestEffortTimestamp(f *Frame) int64 {
	return int64(f.best_effort_timestamp)
}

/** AvGetChannelLayout Return a channel layout id that matches name, or 0 if no match is found.
*
* name can be one or several of the following notations,
* separated by '+' or '|':
* - the name of an usual channel layout (mono, stereo, 4.0, quad, 5.0,
*   5.0(side), 5.1, 5.1(side), 7.1, 7.1(wide), downmix);
* - the name of a single channel (FL, FR, FC, LFE, BL, BR, FLC, FRC, BC,
*   SL, SR, TC, TFL, TFC, TFR, TBL, TBC, TBR, DL, DR);
* - a number of channels, in decimal, followed by 'c', yielding
*   the default channel layout for that number of channels (@see
*   av_get_default_channel_layout);
* - a channel layout mask, in hexadecimal starting with "0x" (see the
*   AV_CH_* macros).
*
* Example: "stereo+FC"  "2c+FC"  "2c+1c"  "0x7"
 */
func AvGetChannelLayout(fmt string) int {
	fmtC := C.CString(fmt)
	defer C.free(unsafe.Pointer(fmtC))
	return int(C.av_get_channel_layout(fmtC))
}

/** AvGetSampleFormat gets sample format from string
* AV_SAMPLE_FMT_U8 u8
* AV_SAMPLE_FMT_S16 s16
* AV_SAMPLE_FMT_S32 s32
* AV_SAMPLE_FMT_FLT flt
* AV_SAMPLE_FMT_DBL dbl
* AV_SAMPLE_FMT_U8P u8p
* AV_SAMPLE_FMT_S16P s16p
* AV_SAMPLE_FMT_S32P s32p
* AV_SAMPLE_FMT_FLTP fltp
* AV_SAMPLE_FMT_DBLP dblp
 */
func AvGetSampleFormat(fmt string) int {
	fmtC := C.CString(fmt)
	defer C.free(unsafe.Pointer(fmtC))
	return int(C.av_get_sample_fmt(fmtC))
}

// GetFrameAudioInfo exports audio frame information from AvFrame
func GetFrameAudioInfo(f *Frame) (af AudioFrame) {
	af.Samples = int64(f.nb_samples)
	af.SampleRate = int64(f.sample_rate)
	af.Pts = int64(f.pts)
	af.ChannelLayout = int64(f.channel_layout)
	af.Format = SampleFormat(f.format)

	return
}

func GetFrameNBSample(f *Frame) int64 {
	return int64(f.nb_samples)
}

// SetFrameAudioInfo Sets audio frame information to avframe
func SetFrameAudioInfo(af AudioFrame, f *Frame) {
	f.nb_samples = C.int(af.Samples)
	f.sample_rate = C.int(af.SampleRate)
	f.pts = C.int64_t(af.Pts)
	f.channel_layout = C.uint64_t(af.ChannelLayout)
	f.format = C.int(af.Format)
}

func (f *Frame) Data() (data [8]*uint8) {
	for i := range data {
		data[i] = (*uint8)(f.data[i])
	}
	return
}

func (f *Frame) Linesize() (linesize [8]int32) {
	for i := range linesize {
		linesize[i] = int32(f.linesize[i])
	}
	return
}

func (f *Frame) Pts() int64 {
	return int64(f.pts)
}
