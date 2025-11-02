package rss

type RssConfig struct {
	ID    string `json:"id"`
	Route string `json:"route"`
}

var RssConfigs = []RssConfig{
	{
		ID:    "Github_trending_typescript",
		Route: "/github/trending/daily/typescript/en",
	},
	{
		ID:    "Github_trending_go",
		Route: "/github/trending/daily/go/en",
	},
	{
		ID:    "Github_trending_javascript",
		Route: "/github/trending/daily/javascript/en",
	},
}