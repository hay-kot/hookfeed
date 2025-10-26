package dtos

type Feed struct {
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	ID          string    `json:"id"`
	Keys        []string  `json:"-"`
	Description string    `json:"description"`
	Middleware  []string  `json:"middleware"`
	Adapters    []string  `json:"adapters"`
	Retention   Retention `json:"retention"`
}

type Retention struct {
	MaxCount   int `json:"maxCount"`
	MaxAgeDays int `json:"maxAgeDays"`
}
