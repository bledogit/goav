package main

import (
	"log"

	"gitlab.com/nielsen-media/eng/reference/commons/goav/avcodec"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avdevice"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avfilter"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avutil"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/swresample"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/swscale"
)

func main() {

	log.Printf("AvFilter Version:\t%v", avfilter.AvfilterVersion())
	log.Printf("AvDevice Version:\t%v", avdevice.AvdeviceVersion())
	log.Printf("SWScale Version:\t%v", swscale.SwscaleVersion())
	log.Printf("AvUtil Version:\t%v", avutil.AvutilVersion())
	log.Printf("AvCodec Version:\t%v", avcodec.AvcodecVersion())
	log.Printf("Resample Version:\t%v", swresample.SwresampleLicense())

}
