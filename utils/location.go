package utils

import (
	"github.com/ip2location/ip2location-go/v9"
)

func GetArea(ip string) (string, error) {
	db, err := ip2location.OpenDB("./conf/ip2location_region.bin")
	if ip == "127.0.0.1" || ip == "::1" {
		return "Local City, Local Country", nil
	}

	if err != nil {
		return "", nil
	}
	results, err := db.Get_all(ip)
	if err != nil {
		return "", nil
	}
	// fmt.Printf("region: %s\n", results.Region)
	return results.Country_long + ", " + results.City, nil
}
