package wrapper

import "C"
import (
	"bufio"
	"fmt"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avcodec"
	"gitlab.com/nielsen-media/eng/reference/commons/goav/avformat"
	"os"
)

type Demuxer struct {
	context  *avformat.Context
	pipefile string
	data     chan []byte
	pipe     *os.File
}

func NewDemuxer() *Demuxer {
	return &Demuxer{}
}

func (d *Demuxer) WithFile(file string) error {
	fmt.Println(avformat.AvformatConfiguration())

	var dacontext *avformat.Context

	dacontext = avformat.AvformatAllocContext()

	if avformat.AvformatOpenInput(&dacontext, file, nil, nil) != 0 {
		return fmt.Errorf("unable to input")
	}

	// Retrieve stream information
	if dacontext.AvformatFindStreamInfo(nil) < 0 {
		return fmt.Errorf("couldn't find stream information")
	}

	dacontext.AvDumpFormat(0, file, 0)

	d.context = dacontext
	return nil
}

func (d *Demuxer) WithPipe(reader *bufio.Reader) error {

	d.pipefile = "/tmp/pipe"

	var err error
	if d.pipe, err = os.OpenFile(d.pipefile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777); err != nil {
		return err
	}

	d.data = make(chan []byte)
	go func() {
		defer d.pipe.Close()
		buffer := make([]byte, 188*10)
		for {
			n, err := reader.Read(buffer)
			if err != nil {
				fmt.Println("break demux reader")
				break
			}
			d.pipe.Write(buffer[:n])
		}
		fmt.Print("Demuxer quit")
	}()

	var dacontext *avformat.Context

	dacontext = avformat.AvformatAllocContext()

	if avformat.AvformatOpenInput(&dacontext, d.pipefile, nil, nil) != 0 {
		return fmt.Errorf("unable to input")
	}

	// Retrieve stream information
	if dacontext.AvformatFindStreamInfo(nil) < 0 {
		return fmt.Errorf("couldn't find stream information")
	}

	dacontext.AvDumpFormat(0, d.pipefile, 0)

	d.context = dacontext
	return nil
}

func (d *Demuxer) Demux(packet *avcodec.Packet) int {
	return d.context.AvReadFrame(packet)
}

func (d *Demuxer) GetContext() *avformat.Context {
	return d.context
}

func (d *Demuxer) Close() {
	d.context.AvformatCloseInput()
}
