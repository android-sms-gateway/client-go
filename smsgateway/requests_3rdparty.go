package smsgateway

import (
	"net/url"
	"strconv"
	"time"
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

// ToURLValues returns the SendOptions as a URL query string in the form of [url.Values].
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

// ListInboxOptions holds optional filters for listing inbox messages.
type ListInboxOptions struct {
	Type     *IncomingMessageType
	Limit    *int
	Offset   *int
	From     *time.Time
	To       *time.Time
	DeviceID *string
}

// ToURLValues returns the ListInboxOptions as URL query parameters.
func (o ListInboxOptions) ToURLValues() url.Values {
	values := url.Values{}
	if o.Type != nil {
		values.Set("type", string(*o.Type))
	}
	if o.Limit != nil {
		values.Set("limit", strconv.Itoa(*o.Limit))
	}
	if o.Offset != nil {
		values.Set("offset", strconv.Itoa(*o.Offset))
	}
	if o.From != nil {
		values.Set("from", o.From.Format(time.RFC3339))
	}
	if o.To != nil {
		values.Set("to", o.To.Format(time.RFC3339))
	}
	if o.DeviceID != nil {
		values.Set("deviceId", *o.DeviceID)
	}
	return values
}

// ListMessagesOptions holds optional filters for listing messages.
type ListMessagesOptions struct {
	From           *time.Time
	To             *time.Time
	State          *string
	DeviceID       *string
	Limit          *int
	Offset         *int
	IncludeContent *bool
}

// ToURLValues returns the ListMessagesOptions as URL query parameters.
func (o ListMessagesOptions) ToURLValues() url.Values {
	values := url.Values{}
	if o.From != nil {
		values.Set("from", o.From.Format(time.RFC3339))
	}
	if o.To != nil {
		values.Set("to", o.To.Format(time.RFC3339))
	}
	if o.State != nil {
		values.Set("state", *o.State)
	}
	if o.DeviceID != nil {
		values.Set("deviceId", *o.DeviceID)
	}
	if o.Limit != nil {
		values.Set("limit", strconv.Itoa(*o.Limit))
	}
	if o.Offset != nil {
		values.Set("offset", strconv.Itoa(*o.Offset))
	}
	if o.IncludeContent != nil {
		values.Set("includeContent", strconv.FormatBool(*o.IncludeContent))
	}
	return values
}
