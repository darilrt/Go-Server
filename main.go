package main

import (
    "go-server/pkg/app"
    "flag"
    "log"
    "net"
)

func main() {
    var (
        host = flag.String("h", "", "Host http address to listen on")
        port = flag.String("p", "8080", "Port nunber for http listener")
    )
    
    addr := net.JoinHostPort(*host, *port)
    err := app.RunHttp(addr)
    
    if err != nil {
        log.Fatal(err);
    }
}