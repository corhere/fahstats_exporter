package main

import (
	"context"
	"strconv"
	"time"

	"github.com/corhere/fahstats_exporter/fahstats"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		"fahstats_scrape_duration_seconds",
		"Total time to scrape Folding@home donor stats",
		nil, nil,
	)

	donorLabels     = []string{"user_name", "user_id"}
	donorCreditDesc = prometheus.NewDesc(
		"fahstats_donor_credit_total",
		"Total Folding@home points earned",
		donorLabels, nil,
	)
	donorWUsDesc = prometheus.NewDesc(
		"fahstats_donor_wus_total",
		"Total Folding@home work units completed",
		donorLabels, nil,
	)
	donorLastWUDesc = prometheus.NewDesc(
		"fahstats_donor_last_wu_seconds",
		"Time elapsed since the last Folding@home Work Unit",
		donorLabels, nil,
	)
	donorActive7dDesc = prometheus.NewDesc(
		"fahstats_donor_7d_active_clients",
		"Active Folding@home clients (within 7 days)",
		donorLabels, nil,
	)
	donorActive50dDesc = prometheus.NewDesc(
		"fahstats_donor_50d_active_clients",
		"Active Folding@home clients (within 50 days)",
		donorLabels, nil,
	)
	donorRankDesc = prometheus.NewDesc(
		"fahstats_donor_rank",
		"Overall Folding@home rank (if points are combined)",
		donorLabels, nil,
	)
	donorTotalUsersDesc = prometheus.NewDesc(
		"fahstats_donor_users_total",
		"Total number of Folding@home users",
		nil, nil,
	)

	teamLabels     = append(donorLabels, "team_name", "team_id")
	teamCreditDesc = prometheus.NewDesc(
		"fahstats_donor_team_credit_total",
		"Total Folding@home points earned for the team",
		teamLabels, nil,
	)
	teamWUsDesc = prometheus.NewDesc(
		"fahstats_donor_team_wus_total",
		"Total Folding@home work units completed for the team",
		teamLabels, nil,
	)
	teamLastWUDesc = prometheus.NewDesc(
		"fahstats_donor_team_last_wu_seconds",
		"Time elapsed since the last Folding@home Work Unit completed for the team",
		teamLabels, nil,
	)
	teamActive7dDesc = prometheus.NewDesc(
		"fahstats_donor_team_7d_active_clients",
		"Active Folding@home clients (within 7 days) contributed for the team",
		teamLabels, nil,
	)
	teamActive50dDesc = prometheus.NewDesc(
		"fahstats_donor_team_50d_active_clients",
		"Active Folding@home clients (within 50 days) contributed for the team",
		teamLabels, nil,
	)
)

// StatsCollector is a custom Prometheus Collector to export
// Folding@Home user statistics.
type StatsCollector struct {
	ctx   context.Context
	donor string
}

// Describe implements the prometheus.Collector interface.
func (c StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

// Collect implements the prometheus.Collector interface.
func (c StatsCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	stats, err := fahstats.Fetch(c.ctx, c.donor)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("fahstats_error", "Error scraping target", nil, nil), err)
		return
	}

	dlabels := []string{stats.Name, strconv.Itoa(stats.ID)}
	ch <- prometheus.MustNewConstMetric(donorCreditDesc, prometheus.UntypedValue, float64(stats.Credit), dlabels...)
	ch <- prometheus.MustNewConstMetric(donorWUsDesc, prometheus.UntypedValue, float64(stats.TotalWUs), dlabels...)
	ch <- prometheus.MustNewConstMetric(donorLastWUDesc, prometheus.GaugeValue, time.Since(stats.LastWorkUnit.Time()).Seconds(), dlabels...)
	ch <- prometheus.MustNewConstMetric(donorActive7dDesc, prometheus.GaugeValue, float64(stats.ActiveClients7), dlabels...)
	ch <- prometheus.MustNewConstMetric(donorActive50dDesc, prometheus.GaugeValue, float64(stats.ActiveClients50), dlabels...)
	ch <- prometheus.MustNewConstMetric(donorRankDesc, prometheus.GaugeValue, float64(stats.Rank), dlabels...)
	ch <- prometheus.MustNewConstMetric(donorTotalUsersDesc, prometheus.GaugeValue, float64(stats.TotalUsers))

	for _, team := range stats.Teams {
		// TODO: assert team.UID == stats.ID?
		tlabels := append(dlabels, team.Name, strconv.Itoa(team.Team))
		ch <- prometheus.MustNewConstMetric(teamCreditDesc, prometheus.UntypedValue, float64(team.Credit), tlabels...)
		ch <- prometheus.MustNewConstMetric(teamWUsDesc, prometheus.UntypedValue, float64(team.TotalWUs), tlabels...)
		ch <- prometheus.MustNewConstMetric(teamLastWUDesc, prometheus.GaugeValue, time.Since(team.LastWorkUnit.Time()).Seconds(), tlabels...)
		ch <- prometheus.MustNewConstMetric(teamActive7dDesc, prometheus.GaugeValue, float64(team.ActiveClients7), tlabels...)
		ch <- prometheus.MustNewConstMetric(teamActive50dDesc, prometheus.GaugeValue, float64(team.ActiveClients50), tlabels...)
	}

	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds())
}
