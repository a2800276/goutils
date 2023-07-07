package goutils

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

type SystemInfo struct {
	osName     string
	osVersion  string
	osArch     string
	osCores    int
	osUser     string
	osHome     string
	osHostname string
	executable string
	goVersion  string
	ifaces     []net.Interface
}

func NewSystemInfo() *SystemInfo {
	// populate the struct
	systemInfo := &SystemInfo{}
	systemInfo.osName = runtime.GOOS
	systemInfo.osVersion = runtime.GOARCH
	systemInfo.osArch = runtime.GOARCH
	systemInfo.osCores = runtime.NumCPU()
	systemInfo.osUser = os.Getenv("USER")
	systemInfo.osHome = os.Getenv("HOME")
	systemInfo.osHostname, _ = os.Hostname()
	systemInfo.executable, _ = os.Executable()
	systemInfo.goVersion = runtime.Version()
	systemInfo.ifaces, _ = net.Interfaces()

	return systemInfo
}

func (s *SystemInfo) String() string {
	str := "OS Name    : " + s.osName + "\n"
	str += "OS Version : " + s.osVersion + "\n"
	str += "OS Arch    : " + s.osArch + "\n"
	str += fmt.Sprintf("OS Cores   : %d\n", s.osCores)
	str += "OS User    : " + s.osUser + "\n"
	str += "OS Home    : " + s.osHome + "\n"
	str += "OS Hostname: " + s.osHostname + "\n"
	str += "Executable : " + s.executable + "\n"
	str += "Go Version : " + s.goVersion + "\n"

	str += "\n\n"
	for _, iface := range s.ifaces {
		str += fmt.Sprintf("Interface: %v\n", iface)
		addrs, err := iface.Addrs()
		if err != nil {
			str += fmt.Sprintf("    Error: %v\n", err)
		} else {
			for _, addr := range addrs {
				str += fmt.Sprintf("     Addr: %v\n", addr)
			}
		}
	}

	str += "\n\n"
	for _, env := range os.Environ() {
		str += fmt.Sprintf("Env: %v\n", env)
	}

	//memStats := runtime.MemStats{}
	//runtime.ReadMemStats(&memStats)
	//str += fmt.Sprintf("Alloc: %v\n", memStats)
	return str
}
