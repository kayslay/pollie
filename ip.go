// Package ipstack provides info on IP address location
// using the http://api.ipstack.com service.

package pollie

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

const (
	MongoIPCollection = "ips"
)

var ipstackURI = "http://api.ipstack.com"

// IPInterfacer describe methods needed to store an IP
type IPInterfacer interface {
	// Get the IP info for a given ip
	Get(ip string) (IPInfo, error)
	Set(IPInfo) error
}

// IPInfo wraps json response
type IPInfo struct {
	IP            string  `json:"ip,omitempty," bson:"ip"`
	Type          string  `json:"type,omitempty," bson:"type"`
	ContinentCode string  `json:"continent_code,omitempty," bson:"continent_code"`
	ContinentName string  `json:"continent_name,omitempty," bson:"continent_name"`
	CountryCode   string  `json:"country_code,omitempty," bson:"country_code"`
	CountryName   string  `json:"country_name,omitempty," bson:"country_name"`
	RegionCode    string  `json:"region_code,omitempty," bson:"region_code"`
	RegionName    string  `json:"region_name,omitempty," bson:"region_name"`
	City          string  `json:"city,omitempty," bson:"city"`
	Zip           string  `json:"zip,omitempty," bson:"zip"`
	Latitude      float64 `json:"latitude,omitempty," bson:"latitude"`
	Longitude     float64 `json:"longitude,omitempty," bson:"longitude"`
	Location      struct {
		GeonameID float64 `json:"geoname_id,omitempty," bson:"geoname_id"`
		Capital   string  `json:"capital,omitempty," bson:"capital"`
		Languages []struct {
			Code   string `json:"code,omitempty," bson:"code"`
			Name   string `json:"name,omitempty," bson:"name"`
			Native string `json:"native,omitempty," bson:"native"`
		} `json:"languages,omitempty," bson:"languages"`
		CountryFlag             string `json:"country_flag,omitempty," bson:"country_flag"`
		CountryFlagEmoji        string `json:"country_flag_emoji,omitempty," bson:"country_flag_emoji"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode,omitempty," bson:"country_flag_emoji_unicode"`
		CallingCode             string `json:"calling_code,omitempty," bson:"calling_code"`
		IsEu                    bool   `json:"is_eu,omitempty," bson:"is_eu"`
	} `json:"location,omitempty," bson:"location"`
}

// MyIP provides information about the public IP address of the client.
func MyIP() (*IPInfo, error) {
	return getInfo(fmt.Sprintf("%s/json", ipstackURI))
}

// ForeignIP provides information about the given IP address (IPv4 or IPv6)
func ForeignIP(ip string, i ...IPInterfacer) (*IPInfo, error) {
	if ip == "" {
		return nil, fmt.Errorf("empty ip address")
	}

	// check if an IPInterfacer is passed
	if len(i) > 0 && i[0] != nil {
		ipInfo, err := i[0].Get(ip)
		// if successfully gotten from the store return
		if err == nil {
			return &ipInfo, nil
		}
	}

	ipInfo, err := getInfo(fmt.Sprintf("%s/%s?access_key=%s", ipstackURI, ip, viper.GetString("IP_STACK")))
	if err != nil {
		return ipInfo, err
	}

	if len(i) > 0 && i[0] != nil {
		i[0].Set(*ipInfo)
	}

	return ipInfo, nil
}

// Undercover code that makes the real call to the webservice
func getInfo(url string) (*IPInfo, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var ipinfo IPInfo
	err = json.NewDecoder(response.Body).Decode(&ipinfo)
	if err != nil {
		return nil, err
	}

	return &ipinfo, nil
}
