package userstack

import "fmt"

type EntityType int

const (
	UnknownEntity EntityType = iota
	Browser
	MobileBrowser
	EmailClient
	App
	FeedReader
	Crawler
	OfflineBrowser
)

var entities = []string{
	"unknown",
	"browser",
	"mobile-browser",
	"email-client",
	"app",
	"feed-reader",
	"crawler",
	"offline-browser",
}

func (e EntityType) String() string {
	return entities[e]
}

// MarshalText satisfies TextMarshaler
func (e EntityType) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText satisfies TextUnmarshaler
func (e *EntityType) UnmarshalText(text []byte) error {
	enum := string(text)
	for i := 0; i < len(entities); i++ {
		if enum == entities[i] {
			*e = EntityType(i)
			return nil
		}
	}
	return fmt.Errorf("unknown userstack entity type: %s", enum)
}

type DeviceType int

const (
	UnknownDevice DeviceType = iota
	Desktop
	Tablet
	Smartphone
	Console
	Smarttv
	Wearable
)

var devices = []string{
	"unknown",
	"desktop",
	"tablet",
	"smartphone",
	"console",
	"smarttv",
	"wearable",
}

func (d DeviceType) String() string {
	return devices[d]
}

// MarshalText satisfies TextMarshaler
func (d DeviceType) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText satisfies TextUnmarshaler
func (d *DeviceType) UnmarshalText(text []byte) error {
	enum := string(text)
	for i := 0; i < len(devices); i++ {
		if enum == devices[i] {
			*d = DeviceType(i)
			return nil
		}
	}
	return fmt.Errorf("unknown userstack device type: %s", enum)
}

type CategoryType int

const (
	UnknownCategory CategoryType = iota
	SearchEngine
	Monitoring
	ScreenshotService
	Scraper
	SecurityScanner
)

var categories = []string{
	"unknown",
	"search-engine",
	"monitoring",
	"screenshot-service",
	"scraper",
	"security-scanner",
}

func (c CategoryType) String() string {
	return categories[c]
}

// MarshalText satisfies TextMarshaler
func (c CategoryType) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText satisfies TextUnmarshaler
func (c *CategoryType) UnmarshalText(text []byte) error {
	enum := string(text)
	for i := 0; i < len(categories); i++ {
		if enum == categories[i] {
			*c = CategoryType(i)
			return nil
		}
	}
	return fmt.Errorf("unknown userstack category type: %s", enum)
}
