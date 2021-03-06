// +build darwin freebsd netbsd openbsd

package pcap

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	device = "/dev/bpf0"
)

type reader struct {
	fd     int
	buflen int // buffer size supplied by bpf
}

type Capture struct {
	header  syscall.BpfHdr
	payload []byte
}

func (r *reader) ReadPacket() (*Capture, error) {
	buf := make([]byte, r.buflen)
	n, e := syscall.Read(r.fd, buf)
	if e != 0 {
		return nil, &os.PathError{"read", device, os.Errno(e)}
	}
	buf = buf[:n]
	header := *(*syscall.BpfHdr)(unsafe.Pointer(&buf[0]))
	capture := &Capture{
		header:  header,
		payload: buf[header.Hdrlen : uint32(header.Hdrlen)+header.Caplen],
	}
	return capture, nil
}

func (r *reader) Close() error {
	syscall.Close(r.fd)
	return nil // TODO(dfc)
}

func ioctl(fd int, request, argp uintptr) error {
	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, argp)
	return os.NewSyscallError("ioctl", int(errorp))
}

func Open() (PacketReader, error) {
	fd, e := syscall.Open(device, os.O_RDONLY|syscall.O_CLOEXEC, 0666)
	if e != 0 {
		return nil, &os.PathError{"open", device, os.Errno(e)}
	}
	var data [16]byte
	data[0] = 'e'
	data[1] = 'n'
	data[2] = '0'

	var len uint32
	var immediate uint32 = 1
	var promisc uint32 = 1
	if err := ioctl(fd, syscall.BIOCGBLEN, uintptr(unsafe.Pointer(&len))); err != nil {
		return nil, err
	}
	if err := ioctl(fd, syscall.BIOCSBLEN, uintptr(unsafe.Pointer(&len))); err != nil {
		return nil, err
	}
	if err := ioctl(fd, syscall.BIOCIMMEDIATE, uintptr(unsafe.Pointer(&immediate))); err != nil {
		return nil, err
	}
	if err := ioctl(fd, syscall.BIOCSETIF, uintptr(unsafe.Pointer(&data[0]))); err != nil {
		return nil, err
	}
	if err := ioctl(fd, syscall.BIOCPROMISC, uintptr(unsafe.Pointer(&promisc))); err != nil {
		return nil, err
	}
	return &reader{fd, int(len)}, nil
}
