// +build linux

package usbserial

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	ErrIsClosed = errors.New("usbserial: is closed")
)

type UsbSerial struct {
	fd     int
	isOpen bool
}

func Open(device string) (*UsbSerial, error) {
	fd, err := unix.Open(device, unix.O_RDWR|unix.O_NOCTTY, 0600)
	if err != nil {
		return nil, err
	}

	rv := &UsbSerial{
		fd: fd,
	}

	cfg := &unix.Termios{}
	if err := rv.ioctl(unix.TCGETS2, uintptr(unsafe.Pointer(cfg))); err != nil {
		unix.Close(fd)
		return nil, err
	}

	cfg.Iflag = 0
	cfg.Oflag = 0
	cfg.Cflag = unix.B115200 | unix.CS8 | unix.CLOCAL | unix.CREAD
	cfg.Lflag = 0
	cfg.Cc[unix.VTIME] = 0
	cfg.Cc[unix.VMIN] = 0

	if err := rv.ioctl(unix.TCSETS2, uintptr(unsafe.Pointer(cfg))); err != nil {
		unix.Close(fd)
		return nil, err
	}

	rv.isOpen = true
	return rv, nil
}

func (u *UsbSerial) ioctl(req uint, arg uintptr) error {
	for {
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(u.fd), uintptr(req), arg)
		if errno == 0 {
			return nil
		}
		if errno != syscall.EINTR {
			return fmt.Errorf("usbserial: %s", errno)
		}
	}
}

func (u *UsbSerial) Close() error {
	if !u.isOpen {
		return nil
	}

	return unix.Close(u.fd)
}

func (u *UsbSerial) Flush() error {
	if !u.isOpen {
		return ErrIsClosed
	}

	return u.ioctl(unix.TCFLSH, uintptr(unix.TCIOFLUSH))
}

func (u *UsbSerial) readByte() (byte, error) {
	if !u.isOpen {
		return 0, ErrIsClosed
	}

	fds := &unix.FdSet{}
	fds.Zero()
	fds.Set(u.fd)

	for {
		ns, err := unix.Select(u.fd+1, fds, nil, nil, nil)
		if err != nil && err != unix.EINTR {
			return 0, err
		}
		if ns == 1 && fds.IsSet(u.fd) {
			break
		}
	}

	b := make([]byte, 1)
	for {
		c, err := unix.Read(u.fd, b)
		if err != nil {
			return 0, err
		}
		if c == 1 {
			return b[0], nil
		}
	}
}

func (u *UsbSerial) ReadLine() (string, error) {
	buf := []byte{}

	for {
		b, err := u.readByte()
		if err != nil {
			return "", err
		}

		buf = append(buf, b)

		if b == '\n' {
			break
		}
	}

	return strings.TrimSpace(string(buf)), nil
}

func (u *UsbSerial) WriteLine(l string, nl bool) error {
	if !u.isOpen {
		return ErrIsClosed
	}

	l = strings.TrimSpace(l)
	if strings.ContainsAny(l, "\r\n") {
		return fmt.Errorf("usbserial: trying to write multiple lines at once: %s", l)
	}

	p := []byte(l)
	if nl {
		p = append(p, '\n')
	}

	n := 0
	for n < len(p) {
		c, err := unix.Write(u.fd, p[n:])
		if err != nil {
			return err
		}
		if c == 0 {
			break
		}
		n += c
	}

	if n != len(p) {
		return fmt.Errorf("usbserial: failed to write full line (%s): %d/%d", l, n, len(p))
	}

	return nil
}
