package models

import "strings"

// Team maps directly to the teams database table.
type Team struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	URLSlug  string `json:"url_slug" db:"url_slug"`
	Homepage string `json:"homepage" db:"homepage"`
	IconURL  string `json:"icon_url" db:"icon_url"`

	LeadID int64 `json:"lead_id" db:"lead_id"`
}

// TeamJSON has additional fields we will be sending to the client.
type TeamJSON struct {
	Team

	Lead User `json:"lead"`
}

// NewTeam makes a new team generating the slug based on it's name.
func NewTeam(name, homepage, iconURL string) Team {
	slug := strings.ToLower(name)
	slug = strings.Replace(slug, " ", "-", -1)

	return Team{
		Name:     name,
		IconURL:  iconURL,
		Homepage: homepage,
		URLSlug:  slug,
	}
}
