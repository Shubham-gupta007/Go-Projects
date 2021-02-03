package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type IP struct {
	IPadddr string `json:"ipaddress"`
}
type Response struct {
	Success bool `json:"success"`
	Data    struct {
		Privateip []string `json:"privateip"`
		Publicip  []string `json:"publicip"`
		Invalidip []string `json:"invalidip"`
	} `json:"data"`
}

func IsPublicIP(IP net.IP) string {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return "private ip"
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return "private ip"
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return "private ip"
		case ip4[0] == 192 && ip4[1] == 168:
			return "private ip"
		default:
			return "public ip"
		}
	}
	return "Invalid IP"
}
func checkip(w http.ResponseWriter, req *http.Request) {
	var ip IP
	var response Response
	var Private []string
	var Public []string
	var Invalid []string

	_ = json.NewDecoder(req.Body).Decode(&ip)

	// var ips = []string{"192.168.0.5", "106.51.36.237", "127.0.0.1", "172.16.25.0", "10.10.25.43", "99.99.99.9999999"}
	address := ip.IPadddr
	fmt.Println(address)

	ips := strings.Split(address, ",")

	for _, ip := range ips {
		IP := net.ParseIP(strings.Trim(ip, " "))

		checkip := IsPublicIP(IP)
		if checkip == "private ip" {
			Private = append(Private, ip)
		} else if checkip == "public ip" {
			Public = append(Public, ip)
		} else if checkip == "Invalid IP" {
			Invalid = append(Invalid, ip)
		}
	}
	response.Success = true
	response.Data.Privateip = Private
	response.Data.Publicip = Public
	response.Data.Invalidip = Invalid
	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	fmt.Println("Starting the application...")
	router.HandleFunc("/iptest", checkip).Methods("POST")
	log.Fatal(http.ListenAndServe(":56789", router))

}
