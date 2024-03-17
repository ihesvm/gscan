package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type ScanResult struct {
	Port   string
	Status string
	Host   string
}

func scanPort(protocol, hostname string, port int) ScanResult {
	result := ScanResult{Port: strconv.Itoa(port) + string("/") + protocol}
	address := hostname + ":" + strconv.Itoa(port)
	ips, _ := net.LookupIP(hostname)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			result.Host = ip.String()
		} else {
			result.Host = hostname
		}
	}
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)

	if err != nil {
		result.Status = "Close"
		return result
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic("Error!")
		}
	}(conn)
	result.Status = "Open"
	return result
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	name := "Object-" + strconv.Itoa(rand.Intn(100)) + ".json"
	_ = os.WriteFile(name, b, 0644)
	return
}

func main() {
	var results []ScanResult

	var port int
	var proto string
	var host string
	var err error

	if len(os.Args) > 1 {
		proto = os.Args[1]
		host = os.Args[2]
		port, err = strconv.Atoi(os.Args[3])
		if err != nil {
			panic(err)
		}
		for i := 0; i < port; i++ {
			results = append(results, scanPort(proto, host, i))
		}

		PrettyPrint(results)

		//fmt.Printf("%v",results)
	} else {
		fmt.Println("Args not found ...")
	}

}
