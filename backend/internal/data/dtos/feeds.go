package dtos

type Feed struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Keys        []string
	Description string    `json:"description"`
	Middleware  []string  `json:"middleware"`
	Adapters    []string  `json:"adapters"`
	Retention   Retention `json:"retention"`
}

type Retention struct {
	MaxCount   int `json:"maxCount"`
	MaxAgeDays int `json:"maxAgeDays"`
}
