package smsgateway

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"
)

type SendOption func(*SendOptions)

type SendOptions struct {
	SkipPhoneValidation *bool `query:"skipPhoneValidation"`
	DeviceActiveWithin  *int  `query:"deviceActiveWithin"  validate:"omitzero,min=1"`
}

func (o *SendOptions) Apply(options ...SendOption) *SendOptions {
	for _, option := range options {
		option(o)
	}

	return o
}

// Validate checks if the SendOptions are valid.
func (o *SendOptions) Validate() error {
	if o.DeviceActiveWithin != nil && (*o.DeviceActiveWithin < 1 || *o.DeviceActiveWithin > math.MaxInt32) {
		return fmt.Errorf("%w: `deviceActiveWithin` must be between 1 and %d", ErrValidationFailed, math.MaxInt32)
	}

	return nil
}

// ToURLValues returns the SendOptions as a URL query string in the form of [url.Values].
// It includes only the options that have been set.
func (o *SendOptions) ToURLValues() url.Values {
	values := url.Values{}
	if o.SkipPhoneValidation != nil {
		values.Set("skipPhoneValidation", strconv.FormatBool(*o.SkipPhoneValidation))
	}
	if o.DeviceActiveWithin != nil {
		values.Set("deviceActiveWithin", strconv.Itoa(*o.DeviceActiveWithin))
	}
	return values
}

// WithSkipPhoneValidation returns a SendOption that disables phone number
// validation for messages. Validation is enabled by default.
func WithSkipPhoneValidation(skipPhoneValidation bool) SendOption {
	return func(o *SendOptions) {
		o.SkipPhoneValidation = &skipPhoneValidation
	}
}

// WithDeviceActiveWithin returns a SendOption that filters devices that have
// been active within the given number of hours.
func WithDeviceActiveWithin(hours uint) SendOption {
	return func(o *SendOptions) {
		var h int
		if hours > math.MaxInt32 {
			h = math.MaxInt32
		} else {
			h = int(hours)
		}

		o.DeviceActiveWithin = &h
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

// ListMessagesOptions holds optional filters and sorting for listing messages.
// Sorting follows the JSON:API specification (sort parameter).
type ListMessagesOptions struct {
	From           *time.Time         `query:"from"`
	To             *time.Time         `query:"to"`
	State          *string            `query:"state"          validate:"omitempty,oneof=Pending Cancelling Cancelled Processed Sent Delivered Failed"`
	DeviceID       *string            `query:"deviceId"       validate:"omitempty,len=21"`
	Limit          *int               `query:"limit"          validate:"omitempty,min=1,max=100"`
	Offset         *int               `query:"offset"         validate:"omitempty,min=0"`
	IncludeContent *bool              `query:"includeContent"`
	Sort           *MessagesSortOrder `query:"sort"           validate:"omitempty,oneof=created_at -created_at"`
}

// Validate checks if the ListMessagesOptions are valid.
func (o ListMessagesOptions) Validate() error {
	if o.From != nil && o.To != nil && o.From.After(*o.To) {
		return fmt.Errorf("%w: `from` date must be before `to` date", ErrValidationFailed)
	}

	return nil
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
	if o.Sort != nil {
		values.Set("sort", string(*o.Sort))
	}
	return values
}

type MessagesSortOrder string

const (
	CreatedAtAscending  MessagesSortOrder = "created_at"
	CreatedAtDescending MessagesSortOrder = "-created_at"
)
