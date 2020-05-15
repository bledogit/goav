package wrapper

import (
	"bufio"
	"fmt"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"os"
)

type Demuxer struct {
	context  *avformat.Context
	pipefile string
	data     chan []byte
	pipe     *os.File
	packet   *avcodec.Packet
}

func (d *Demuxer) InitWithFile(file string) error {

	d.packet = avcodec.AvPacketAlloc()

	var dacontext *avformat.Context

	dacontext = avformat.AvformatAllocContext()

	if avformat.AvformatOpenInput(&dacontext, file, nil, nil) != 0 {
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

func (d *Demuxer) InitWithPipe(reader *bufio.Reader) error {

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

	d.packet = avcodec.AvPacketAlloc()

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

func (d *Demuxer) Demux() *avcodec.Packet {
	if code := d.context.AvReadFrame(d.packet); code >= 0 {
		return d.packet
	}
	return nil
}

func (d *Demuxer) GetContext() *avformat.Context {
	return d.context
}

func (d *Demuxer) Close() {

}
