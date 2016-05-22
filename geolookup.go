//
// geolookup.go
// Messy code to lookup an ip in a sqlite database
//

package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
)

// Inet_Aton converts an IPv4 net.IP object to a 64 bit integer.
func Inet_Aton(ip net.IP) int64 {
	ipv4Int := big.NewInt(0)
	ipv4Int.SetBytes(ip.To4())
	return ipv4Int.Int64()
}

// Checks ip to see if it falls within reserved ip range
func IsPrivateIP(ip net.IP) bool {
	cidr := []string { "0.0.0.0/8", "10.0.0.0/8", "100.64.0.0/10",
			   "169.254.0.0/16", "172.16.0.0/12", "192.0.0.0/24",
			   "192.0.0.0/24", "192.0.2.0/24", "192.88.99.0/24",
			   "192.168.0.0/16", "198.18.0.0/15", "198.51.100.0/24",
			   "203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4" }

	for i := 0; i < len(cidr); i++ {
		_, ipn,_ := net.ParseCIDR(cidr[i])
		if ipn.Contains(ip) {
			return true
		} 
	}
	return false
}

// TODO
// - Clean up the code
func main() {
	var ip_i int64 = 0
	var banner string = "Geolookup v0.01 by Gau Bac Cuc"
	
	fmt.Println(banner)
	
	if(len(os.Args) < 2) {
		fmt.Println("Must supply an ip address")
		return
	}

	ip_s := os.Args[1]

	ip := net.ParseIP(ip_s)
	if(ip != nil) {
		ip_i = Inet_Aton(ip)
	} else {
		fmt.Println("Invalid ip address")
		return
	}

	// Check the ip to make sure its valid
	if ip.IsLoopback() {
		fmt.Printf("Error: %s is a loopback address\n", ip_s)
		return
	}
	if ip.IsMulticast() {
		fmt.Printf("Error: %s is a multicast address\n", ip_s)
		return
	}
	if IsPrivateIP(ip) {
	        fmt.Printf("Error: %s is a reserved address\n", ip_s)
                return
	}

	// Try to open the database
	db, err :=  sql.Open("sqlite3", "ip2nation.db")
	if(db == nil) {
		fmt.Printf("Unable to open db error %s\n", err)
	} else {
		fmt.Printf("Database - http://www.ip2nation.com\n\n")
	}

	strQuery := fmt.Sprintf("SELECT c.country FROM ip2nationCountries c, ip2nation i WHERE i.ip < %d AND c.code = i.country ORDER BY i.ip DESC LIMIT 0,1;", ip_i)
	
	rows, err := db.Query(strQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	defer rows.Close()
	rows.Next()

	var country string

	err = rows.Scan(&country)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s resolves to: %s\n", ip_s, country)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
