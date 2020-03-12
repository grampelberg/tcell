package tcell

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var ErrWinSizeUnused = errors.New("driver does not provide WinSize")

type TermDriver interface {
	Init(winch chan os.Signal) (*os.File, *os.File, error)
	Fini()
	WinSize() (int, int, error)
}

type defaultTermDriver struct {
	winch chan os.Signal
	out   *os.File
}

func (d *defaultTermDriver) Init(winch chan os.Signal) (in *os.File, out *os.File, err error) {
	in, err = os.OpenFile("/dev/tty", os.O_RDONLY, 0)
	if err != nil {
		return
	}

	out, err = os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		return
	}

	signal.Notify(winch, syscall.SIGWINCH)
	d.winch = winch
	d.out = out
	return
}

func (d *defaultTermDriver) Fini() {
	signal.Stop(d.winch)
}

func (d *defaultTermDriver) WinSize() (int, int, error) {
	return 0, 0, ErrWinSizeUnused
}
