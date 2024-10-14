package model

import (
	"sort"
	"time"
)

type ChatWorkReport struct {
	UserId             string                        `json:"user_id"`
	UserFullname       string                        `json:"user_fullname"`
	Total              int                           `json:"total"`
	Facebook           ChannelWorkPerformanceMetrics `json:"facebook"`
	Zalo               ChannelWorkPerformanceMetrics `json:"zalo"`
	ConversationExists map[string]bool               `json:"-"`
}

type ChannelWorkPerformanceMetrics struct {
	TotalChannels int                `json:"total_channels"`
	ReceivingTime PerformanceMetrics `json:"receiving_time"` // metrics of the time between first time user response
	ReplyingTime  PerformanceMetrics `json:"replying_time"`  // metrics of the time between each time user replies
}

type ChatGeneralReport struct {
	Channel            string          `json:"channel"`
	OaName             string          `json:"oa_name"`
	TotalConversations int             `json:"total_conversations"`
	Fresh              QuantityRatio   `json:"fresh"`
	Processing         QuantityRatio   `json:"processing"`
	Resolved           QuantityRatio   `json:"resolved"`
	ConversationExists map[string]bool `json:"-"`
}

type QuantityRatio struct {
	Quantity int `json:"quantity"`
	Percent  int `json:"percent"`
}

type PerformanceMetrics struct {
	Timestamps []time.Duration `json:"-"`
	Fastest    int             `json:"fastest"`
	Average    int             `json:"average"`
	Slowest    int             `json:"slowest"`
}

func (p *PerformanceMetrics) AddTimestamp(timestamp time.Duration) {
	if p.Timestamps == nil {
		p.Timestamps = make([]time.Duration, 0)
	}
	// timestamp should not be negative numbers
	if timestamp.Milliseconds() < 0 {
		return
	}
	p.Timestamps = append(p.Timestamps, timestamp)
}

func (p *PerformanceMetrics) CalculateMetrics() {
	if len(p.Timestamps) == 0 {
		return
	}

	// sort ascending, fastest -> slowest
	sort.Slice(p.Timestamps, func(i, j int) bool {
		return p.Timestamps[i] < p.Timestamps[j]
	})

	p.Fastest = int(p.Timestamps[0].Milliseconds())
	p.Slowest = int(p.Timestamps[len(p.Timestamps)-1].Milliseconds())
	var totalMs int64
	for _, ts := range p.Timestamps {
		totalMs += ts.Milliseconds()
	}
	p.Average = int(totalMs / int64(len(p.Timestamps)))
}
