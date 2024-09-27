package sender

import (
	v1 "github.com/llmariner/api-usage/api/v1"
)

// UsageSetter is an interface that allows components to add usage data to the sender.
type UsageSetter interface {
	AddUsage(usage *v1.UsageRecord)
}

// NoopUsageSetter does nothing.
type NoopUsageSetter struct{}

// AddUsage does nothing
func (s NoopUsageSetter) AddUsage(usage *v1.UsageRecord) {}
