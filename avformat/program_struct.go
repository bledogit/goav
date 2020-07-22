package avformat

import "C"
import (
	"unsafe"
)

type Program struct {
	Pid           int
	ProgramNumber int
	StreamIndexes []int
	PcrPid        int
	PmtVersion    int
}

func (prg *AvProgram) GetProgram() Program {

	nstreams := int(prg.nb_stream_indexes)

	data := (*[1 << 28]C.uint)(unsafe.Pointer(unsafe.Pointer(prg.stream_index)))[:nstreams:nstreams]

	streams := []int{}
	for _, s := range data {
		streams = append(streams, int(s))
	}

	p := Program{
		Pid:           int(prg.pmt_pid),
		ProgramNumber: int(prg.program_num),
		StreamIndexes: streams,
		PcrPid:        int(prg.pcr_pid),
		PmtVersion:    int(prg.pmt_version),
	}

	return p
}
