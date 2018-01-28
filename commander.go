package sasrd

import (
	"os/exec"
	"runtime"
)

// Commander handles the requests to reboot, shutdown or hibernate the system
type Commander interface {
	Reboot() *exec.Cmd
	Shutdown() *exec.Cmd
	Hibernate() *exec.Cmd
}

func NewCommander() Commander {
	switch os := runtime.GOOS; os {
	case "darwin":
		return new(osxCommander)
	case "linux":
		return new(linuxCommander)
	case "windows":
		return new(winCommander)
	default:
		return nil
	}
}

type winCommander struct {
}

func (o *winCommander) Reboot() *exec.Cmd {
	return exec.Command("shutdown.exe", "-r", "-f")
}

func (o *winCommander) Shutdown() *exec.Cmd {
	return exec.Command("shutdown.exe", "-s", "-f")
}

func (o *winCommander) Hibernate() *exec.Cmd {
	return exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "1")
}

type osxCommander struct {
}

func (o *osxCommander) Reboot() *exec.Cmd {
	return exec.Command("shutdown", "-r", "now")
}

func (o *osxCommander) Shutdown() *exec.Cmd {
	return exec.Command("shutdown", "-h", "now")
}

func (o *osxCommander) Hibernate() *exec.Cmd {
	return exec.Command("shutdown", "-s", "now")
}

type linuxCommander struct {
}

func (o *linuxCommander) Reboot() *exec.Cmd {
	return exec.Command("shutdown", "-r", "now")
}

func (o *linuxCommander) Shutdown() *exec.Cmd {
	return exec.Command("shutdown", "-h", "now")
}

func (o *linuxCommander) Hibernate() *exec.Cmd {
	return exec.Command("pm-hibernate")
}
