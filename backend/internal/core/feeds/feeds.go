package feeds

const (
	DefaultRetentionCount   = 10_000
	DefaultRetentionMaxDays = 10_000
	DefaultAdaptersEnabled  = true
)

// Config represents the complete HookFeed configuration
type Config struct {
	Middleware []string `yaml:"middleware"` // Filenames in execution order
	Feeds      []Feed   `yaml:"feeds"`
}

// Feed represents a webhook feed configuration
type Feed struct {
	Name            string     `yaml:"name"`
	Category        string     `yaml:"category"`
	ID              string     `yaml:"id"`   // used as the unique identifier
	Keys            []string   `yaml:"keys"` // used as the :key value in url path to resolve feed
	Description     string     `yaml:"description"`
	Middleware      []string   `yaml:"middleware"` // Filenames of middleware scripts
	AdaptersEnabled *bool      `yaml:"adapters_enabled"`
	Adapters        []string   `yaml:"adapters"` // pointer to distinguish between null, empty array, and populated array
	Retention       *Retention `yaml:"retention"`
}

func (f Feed) IntoParsed() FeedParsed {
	fp := FeedParsed{
		Name:            f.Name,
		Category:        f.Category,
		ID:              f.ID,
		Keys:            f.Keys,
		Description:     f.Description,
		Middleware:      f.Middleware,
		AdaptersEnabled: DefaultAdaptersEnabled,
		Adapters:        f.Adapters,
		Retention: RetentionParsed{
			MaxCount:   DefaultRetentionCount,
			MaxAgeDays: DefaultRetentionMaxDays,
		},
	}

	if f.AdaptersEnabled != nil {
		fp.AdaptersEnabled = *f.AdaptersEnabled
	}

	if f.Retention != nil {
		if f.Retention.MaxCount != nil {
			fp.Retention.MaxCount = *f.Retention.MaxCount
		}

		if f.Retention.MaxAgeDays != nil {
			fp.Retention.MaxAgeDays = *f.Retention.MaxAgeDays
		}
	}

	return fp
}

// FeedParsed is the valid verion of [Feed] where no properties are unset. This struct has default values
// where none were assigned in the base type
type FeedParsed struct {
	Name            string          `yaml:"name"`
	Category        string          `yaml:"category"`
	ID              string          `yaml:"id"`   // used as the unique identifier
	Keys            []string        `yaml:"keys"` // used as the :key value in url path to resolve feed
	Description     string          `yaml:"description"`
	Middleware      []string        `yaml:"middleware"` // Filenames of middleware scripts
	AdaptersEnabled bool            `yaml:"adapters_enabled"`
	Adapters        []string        `yaml:"adapters"` // pointer to distinguish between null, empty array, and populated array
	Retention       RetentionParsed `yaml:"retention"`
}

// Retention defines message retention policies
type Retention struct {
	MaxCount   *int `yaml:"max_count"`
	MaxAgeDays *int `yaml:"max_age_days"`
}

// RetentionParsed defines message retention policies that have been parsed into a valid
// struct (no optional values)
type RetentionParsed struct {
	MaxCount   int `yaml:"max_count"`
	MaxAgeDays int `yaml:"max_age_days"`
}
