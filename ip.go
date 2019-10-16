// Package ipstack provides info on IP address location
// using the http://api.ipstack.com service.

package pollie

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

var ipstackURI = "http://api.ipstack.com"

// IPInterfacer describe methods needed to store an IP
type IPInterfacer interface {
	Get(ip string) (*IPInfo, error)
	Set(IPInfo) error
}

// IPInfo wraps json response
type IPInfo struct {
	IP            string  `json:"ip,omitempty"`
	Type          string  `json:"type,omitempty"`
	ContinentCode string  `json:"continent_code,omitempty"`
	ContinentName string  `json:"continent_name,omitempty"`
	CountryCode   string  `json:"country_code,omitempty"`
	CountryName   string  `json:"country_name,omitempty"`
	RegionCode    string  `json:"region_code,omitempty"`
	RegionName    string  `json:"region_name,omitempty"`
	City          string  `json:"city,omitempty"`
	Zip           string  `json:"zip,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	Location      struct {
		GeonameID float64 `json:"geoname_id,omitempty"`
		Capital   string  `json:"capital,omitempty"`
		Languages []struct {
			Code   string `json:"code,omitempty"`
			Name   string `json:"name,omitempty"`
			Native string `json:"native,omitempty"`
		} `json:"languages,omitempty"`
		CountryFlag             string `json:"country_flag,omitempty"`
		CountryFlagEmoji        string `json:"country_flag_emoji,omitempty"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode,omitempty"`
		CallingCode             string `json:"calling_code,omitempty"`
		IsEu                    bool   `json:"is_eu,omitempty"`
	} `json:"location,omitempty"`
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
	// checki if an IPInterfacer is passed
	if len(i) > 0 && i[0] != nil {
		ipInfo, err := i[0].Get(ip)
		// if successfully gotten from the store return
		if err == nil {
			return ipInfo, nil
		}
	}

	ipInfo, err := getInfo(fmt.Sprintf("%s/%s?access_key=%s", ipstackURI, ip, viper.GetString("IP_STACK")))
	if len(i) > 0 && i[0] != nil {
		i[0].Set(*ipInfo)
	}

	return ipInfo, err
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
