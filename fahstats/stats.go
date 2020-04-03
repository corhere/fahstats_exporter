package fahstats

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const url = "https://stats.foldingathome.org/api/donor/"

type FahTime time.Time

func (t *FahTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	v, err := time.Parse(`"2006-01-02 15:04:05"`, string(data))
	*t = FahTime(v)
	return err
}

func (t FahTime) Time() time.Time {
	return (time.Time)(t)
}

// base declares all the stats fields common to both users and teams.
type base struct {
	// Total WUs
	TotalWUs     int     `json:"wus"`
	LastWorkUnit FahTime `json:"last"`
	Name         string  `json:"name"`
	Credit       int     `json:"credit"`
	// Active clients (within 50 days)
	ActiveClients50 int `json:"active_50"`
	// Active clients (within 7 days)
	ActiveClients7 int `json:"active_7"`
}

// Team are the statistics for a team.
type Team struct {
	Team int `json:"team"`
	UID  int `json:"uid"`
	base
}

// Donor are the statistics for an individual donor.
type Donor struct {
	ID int `json:"id"`
	// Overall rank (if points are combined)
	Rank int `json:"rank"`
	// Total (active?) users
	TotalUsers int `json:"total_users"`
	base
	Teams []Team `json:"teams"`
}

// Fetch gets the current stats for the specified donor.
func Fetch(ctx context.Context, donor string) (Donor, error) {
	req, err := http.NewRequest("GET", url+donor, nil)
	if err != nil {
		return Donor{}, err
	}
	req.Header.Add("User-Agent", "fahstats_exporter/0.0.1")
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return Donor{}, err
	}

	var stats Donor
	if err = json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return Donor{}, err
	}

	return stats, nil
}
