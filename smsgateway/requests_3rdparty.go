package smsgateway

import (
	"net/url"
	"strconv"
)

type SendOption func(*SendOptions)

type SendOptions struct {
	skipPhoneValidation *bool
	deviceActiveWithin  *uint
}

func (o *SendOptions) Apply(options ...SendOption) *SendOptions {
	for _, option := range options {
		option(o)
	}

	return o
}

// ToURLValues returns the SendOptions as a URL query string in the form of url.Values.
// It includes only the options that have been set.
func (o *SendOptions) ToURLValues() url.Values {
	values := url.Values{}
	if o.skipPhoneValidation != nil {
		values.Set("skipPhoneValidation", strconv.FormatBool(*o.skipPhoneValidation))
	}
	if o.deviceActiveWithin != nil {
		values.Set("deviceActiveWithin", strconv.FormatUint(uint64(*o.deviceActiveWithin), 10))
	}
	return values
}

// WithSkipPhoneValidation returns a SendOption that disables phone number
// validation for messages. Validation is enabled by default.
func WithSkipPhoneValidation(skipPhoneValidation bool) SendOption {
	return func(o *SendOptions) {
		o.skipPhoneValidation = &skipPhoneValidation
	}
}

// WithDeviceActiveWithin returns a SendOption that filters devices that have
// been active within the given number of hours.
func WithDeviceActiveWithin(hours uint) SendOption {
	return func(o *SendOptions) {
		o.deviceActiveWithin = &hours
	}
}
