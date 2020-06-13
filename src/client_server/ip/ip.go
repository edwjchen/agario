package ip

import (
    "net"
    "strings"
)

const INTERFACE_NAME = "eth0"

func GetLocalIP() string {
    l, _ := net.InterfaceByName(INTERFACE_NAME)
    addrs, _ := l.Addrs()
    addr := addrs[0]
    return strings.Split(addr.String(), "/")[0]
}
