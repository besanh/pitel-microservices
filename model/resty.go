package model

import (
	"database/sql"
	"time"
)

type (
	Resty interface{}

	RestyAuth struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	RestySetting struct {
		Url                string        `json:"url"`
		Accept             string        `json:"accept"`
		AuthType           string        `json:"auth_type"`
		RestyAuth          RestyAuth     `json:"resty_auth"`
		Timeout            time.Duration `json:"timeout"`
		InsecureSkipVerify sql.NullBool  `json:"insecure_skip_verify"`
	}

	RestyResponse struct {
		StatusCode int           `json:"status_code"`
		Status     string        `json:"status"`
		Proto      any           `json:"proto"`
		Time       time.Duration `json:"time"`
		ReceivedAt time.Time     `json:"received_at"`
		Body       any           `json:"body"`
	}

	RestyTrace struct {
		DNSLookup      time.Duration `json:"dns_lookup"`
		ConnTime       time.Duration `json:"conn_time"`
		TCPConnTime    time.Duration `json:"tcp_conn_time"`
		TLSHandshake   time.Duration `json:"tls_handshake"`
		ServerTime     time.Duration `json:"server_time"`
		TotalTime      time.Duration `json:"total_time"`
		IsConnReused   bool          `json:"is_conn_reused"`
		IsConnWasIdle  bool          `json:"is_conn_was_idle"`
		ConnIdleTime   time.Duration `json:"conn_idle_time"`
		RequestAttempt int           `json:"request_attempt"`
		RemoteAddr     string        `json:"remote_addr"`
	}
)
