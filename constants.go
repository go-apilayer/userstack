package userstack

type EntityType string

const (
	EntityUnknown        EntityType = "unknown"
	EntityBrowser        EntityType = "browser"
	EntityMobileBrowser  EntityType = "mobile-browser"
	EntityEmailClient    EntityType = "email-client"
	EntityApp            EntityType = "app"
	EntityFeedReader     EntityType = "feed-reader"
	EntityCrawler        EntityType = "crawler"
	EntityOfflineBrowser EntityType = "offline-browser"
)

func (e EntityType) String() string { return string(e) }

// UnmarshalText satisfies TextUnmarshaler
func (e *EntityType) UnmarshalText(text []byte) error {
	enum := string(text)
	if !strictUnmarshal {
		*e = EntityType(enum)
		return nil
	}
	switch enum {
	case "unknown":
		*e = EntityUnknown
	case "browser":
		*e = EntityBrowser
	case "mobile-browser":
		*e = EntityMobileBrowser
	case "email-client":
		*e = EntityEmailClient
	case "app":
		*e = EntityApp
	case "feed-reader":
		*e = EntityFeedReader
	case "crawler":
		*e = EntityCrawler
	case "offline-browser":
		*e = EntityOfflineBrowser
	default:
		return &UnsupportedTypeError{typ: enum, fieldName: "entity"}
	}
	return nil
}

type DeviceType string

const (
	DeviceUnknown    DeviceType = "unknown"
	DeviceDesktop    DeviceType = "desktop"
	DeviceTablet     DeviceType = "tablet"
	DeviceSmartPhone DeviceType = "smartphone"
	DeviceConsole    DeviceType = "console"
	DeviceSmartTV    DeviceType = "smarttv"
	DeviceWearable   DeviceType = "wearable"
)

func (d DeviceType) String() string { return string(d) }

// UnmarshalText satisfies TextUnmarshaler
func (d *DeviceType) UnmarshalText(text []byte) error {
	enum := string(text)
	if !strictUnmarshal {
		*d = DeviceType(enum)
		return nil
	}
	switch enum {
	case "unknown":
		*d = DeviceUnknown
	case "desktop":
		*d = DeviceDesktop
	case "tablet":
		*d = DeviceTablet
	case "smartphone":
		*d = DeviceSmartPhone
	case "console":
		*d = DeviceConsole
	case "smarttv":
		*d = DeviceSmartTV
	case "wearable":
		*d = DeviceWearable
	default:
		return &UnsupportedTypeError{typ: enum, fieldName: "device"}
	}
	return nil
}

type CategoryType string

const (
	CategoryUnknown           CategoryType = "unknown"
	CategorySearchEngine      CategoryType = "search-engine"
	CategoryMonitoring        CategoryType = "monitoring"
	CategoryScreenshotService CategoryType = "screenshot-service"
	CategoryScraper           CategoryType = "scraper"
	CategorySecurityScanner   CategoryType = "security-scanner"
)

func (c CategoryType) String() string { return string(c) }

// UnmarshalText satisfies TextUnmarshaler
func (c *CategoryType) UnmarshalText(text []byte) error {
	enum := string(text)
	if !strictUnmarshal {
		*c = CategoryType(enum)
		return nil
	}
	switch enum {
	case "unknown":
		*c = CategoryUnknown
	case "search-engine":
		*c = CategorySearchEngine
	case "monitoring":
		*c = CategoryMonitoring
	case "screenshot-service":
		*c = CategoryScreenshotService
	case "scraper":
		*c = CategoryScraper
	case "security-scanner":
		*c = CategorySecurityScanner
	default:
		return &UnsupportedTypeError{typ: enum, fieldName: "category"}
	}
	return nil
}
