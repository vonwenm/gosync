package utils

import (
    "github.com/BurntSushi/toml"
    "io/ioutil"
    "net"
    "sync"
    "strings"
)

var (
    config     *Configuration
    configLock = new(sync.RWMutex)
)

func ReadConfigFromFile(configfile string) {
    config_file, err := ioutil.ReadFile(configfile)
    if err != nil {
        panic(err.Error())
    }
    tempConf := new(Configuration)
    _, err = toml.Decode(string(config_file), &tempConf)
    if err != nil {
        panic(err.Error())
    }
    configLock.Lock()
    config = tempConf
    configLock.Unlock()
}

func GetConfig() *Configuration {
    configLock.RLock()
    defer configLock.RUnlock()
    return config
}

func GetLocalIp() string {
    //LogWriteF("Ip address: %v",ip)
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        Check(err, 500, "No network interfaces found")
    }
    var ips []string
    for _, a := range addrs {
        if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                ips = append(ips, ipnet.IP.String())
            }
        }
    }
    return strings.Join(ips, ",")
}
