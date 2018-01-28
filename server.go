package sasrd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/grandcat/zeroconf"
)

type ServerConfig struct {
	Host      string
	Port      int
	Cert      string
	Key       string
	Pass      string
	Interface string
}

type ServiceInfo struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Hostname  string `json:"name"`
	GOOS      string `json:"os"`
	GOARCH    string `json:"arch"`
	SSL       bool   `json:"ssl"`
	Protected bool   `json:"protected"`
}

func NewConfig() ServerConfig {
	return ServerConfig{
		Host: "",
		Port: 38765,
		Cert: "",
		Key:  "",
		Pass: "",
	}
}

func NewConfigFromFile(file string) (ServerConfig, error) {
	result := NewConfig()
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func StartServer(config ServerConfig) {
	interfaces := listInterfaces(func(ifi net.Interface) bool {
		return config.Interface == "" || config.Interface == ifi.Name
	})
	if interfaces == nil || len(interfaces) == 0 {
		logFatal("Could not determine host IP addresses")
	}
	logInfo("Starting Service at interface: %v", interfaces[0].Name)

	hostname, _ := os.Hostname()
	serviceName := fmt.Sprintf("SASR Daemon (%s)", hostname)
	server, err := zeroconf.Register(serviceName, "_http._tcp", "local.", config.Port, []string{
		fmt.Sprintf("host=%v", config.Host),
		fmt.Sprintf("port=%v", config.Port),
		fmt.Sprintf("name=%v", hostname),
		fmt.Sprintf("os=%v", runtime.GOOS),
		fmt.Sprintf("arch=%v", runtime.GOARCH),
		fmt.Sprintf("ssl=%v", config.Key != "" && config.Cert != ""),
		fmt.Sprintf("protected=%v", config.Pass != ""),
		fmt.Sprintf("mac=%v", interfaces[0].HardwareAddr),
	}, nil)
	if err != nil {
		logFatal(err.Error())
	}
	defer server.Shutdown()

	http.HandleFunc("/status", middleware(config, handleStatus))
	http.HandleFunc("/reboot", middleware(config, handleReboot))
	http.HandleFunc("/hibernate", middleware(config, handleHibernate))
	http.HandleFunc("/shutdown", middleware(config, handleShutdown))

	hostWithPort := fmt.Sprintf("%v:%v", config.Host, config.Port)

	if config.Cert != "" && config.Key != "" {
		logInfo("Starting TLS server at %v", hostWithPort)
		err = http.ListenAndServeTLS(hostWithPort, config.Cert, config.Key, nil)
	} else {
		logInfo("Starting server at %v", hostWithPort)
		err = http.ListenAndServe(hostWithPort, nil)
	}
	if err != nil {
		logFatal(err.Error())
	}
}

func middleware(config ServerConfig, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
			}
		}()

		logInfo("%s %s", r.Method, r.URL.Path)

		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		if config.Pass != "" && config.Pass != r.Header.Get("X-PASS") {
			http.Error(w, "Not authenticated", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func checkError(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Status"))
}

func handleReboot(w http.ResponseWriter, r *http.Request) {
	cmd := NewCommander().Reboot()
	out, err := cmd.CombinedOutput()
	checkError(err)
	w.Write(out)
}

func handleHibernate(w http.ResponseWriter, r *http.Request) {
	cmd := NewCommander().Hibernate()
	out, err := cmd.CombinedOutput()
	checkError(err)
	w.Write(out)
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	cmd := NewCommander().Shutdown()
	out, err := cmd.CombinedOutput()
	checkError(err)
	w.Write(out)
}

func listInterfaces(filter ...func(net.Interface) bool) []net.Interface {
	var interfaces []net.Interface
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, ifi := range ifaces {
		if (ifi.Flags & net.FlagUp) == 0 {
			continue
		}
		if (ifi.Flags & net.FlagMulticast) == 0 {
			continue
		}
		if len(filter) == 1 && !filter[0](ifi) {
			continue
		}
		interfaces = append(interfaces, ifi)
	}

	return interfaces
}
