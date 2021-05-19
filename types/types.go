package types

type Site struct {
	Url string `json:"url"`
}

type State struct {
	Blacklist []Site `json:"blacklist"`
}
