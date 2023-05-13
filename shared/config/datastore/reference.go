package datastore

import "time"

const (
	defaultTimeToLiveMs = 900000 // 15 minutes
	defaultRetryTimeMs  = 5000   // 5 seconds
)

//Reference datastore reference
type Reference struct {
	Connection   string
	Namespace    string
	Dataset      string
	TimeToLiveMs int `json:",omitempty" yaml:",omitempty"`
	timeToLive   time.Duration
	RetryTimeMs  int `json:",omitempty" yaml:",omitempty"`
	retryTime    time.Duration
	ReadOnly 