package service

import (
	"log/slog"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPService struct {
	db *geoip2.Reader
}

func NewGeoIPService(dbPath string) *GeoIPService {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		slog.Error("error to load geoip2 service", "error", err)
		os.Exit(1)
	}

	return &GeoIPService{db: db}
}

func (s *GeoIPService) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *GeoIPService) GetCountryCodeFromIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		slog.Debug("failed to parse IP", "ip", ipStr)
		return "UNKNOWN"
	}

	countryCode, err := s.db.Country(ip)
	if err != nil {
		slog.Debug("failed to lookup IP", "ip", ipStr, "error", err)
		return "UNKNOWN"
	}

	if countryCode.Country.IsoCode == "" {
		return "UNKNOWN"
	}

	return countryCode.Country.IsoCode
}
