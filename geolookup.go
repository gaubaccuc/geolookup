// geolookup.go
//
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
)

// Inet_Aton converts an IPv4 net.IP object to a 64 bit integer.
func Inet_Aton(ip net.IP) int64 {
	ipv4Int := big.NewInt(0)
	ipv4Int.SetBytes(ip.To4())
	return ipv4Int.Int64()
}

// TODO
// - Clean up the code
// - Add user input instead of static variable
func main() {
	var ip_i int64 = 0
	ip_s := "8.8.8.8"

	var banner string = "Geolookup v0.01 by Gau Bac Cuc"

	fmt.Println(banner)
		
	ip := net.ParseIP(ip_s)
	if(ip != nil) {
		ip_i = Inet_Aton(ip)
	} else {
		return
	}

	db,err :=  sql.Open("sqlite3", "ip2nation.db")
	if(db == nil) {
		fmt.Printf("Unable to open db error %s\n", err)
	} else {
		fmt.Printf("->DB: http://www.ip2nation.com\n\n")
	}

	strQuery := ""
	
	var country string

	strQuery = fmt.Sprintf("SELECT c.country FROM ip2nationCountries c, ip2nation i WHERE i.ip < %d AND c.code = i.country ORDER BY i.ip DESC LIMIT 0,1;", ip_i)
	
	rows, err := db.Query(strQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		err = rows.Scan(&country)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s resolves to: %s\n",ip_s, country)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
