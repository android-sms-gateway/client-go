# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### New Features
#### Message scheduling
- **Added `Message.ValidUntilTime(now)`** — resolves the effective expiry of a message from `ValidUntil`, or from `TTL` computed against `now`, returning `nil` when neither is set.

### Bug Fixes
- **`Message.Validate()` rejects impossible schedules** — a message whose `scheduleAt` is later than its effective expiry (`validUntil` or `ttl`) now fails validation with an error wrapping `ErrValidationFailed`.

## [1.15.0] - 2026-08-31

### New Features
#### Outgoing MMS
- **Send MMS messages** — `Message.MmsMessage` accepts an optional subject, text, and a list of base64-encoded attachments (`MmsAttachment`). `Message.Validate()` now accepts exactly one of text, data, or MMS content, and requires an MMS to carry either text or at least one attachment.
- **MMS content in message states** — `MessageState.MmsMessage` exposes MMS content when message state is requested with `includeContent=true`.

## [1.14.6] - 2026-08-21

### Bug Fixes
- **Duplicate recipient phone numbers rejected** — a message that lists the same phone number more than once in `PhoneNumbers` now fails validation instead of being accepted.

## [1.14.5] - 2026-08-18

### New Features
#### Batch webhook events
- **Batch event types and payloads** — added `sms:batch:received`, `sms:batch:data-received`, `mms:batch:received`, and `mms:batch:downloaded` event constants with typed payloads carrying the ordered list of messages. Batch events are included in `WebhookEventTypes()`.
#### Inbox refresh
- **Added `WebhookDelivery` to inbox refresh** — `InboxRefreshRequest.WebhookDelivery` selects `Disabled`, `Individual` (one webhook per message), or `Batch` webhook delivery when refreshing the inbox. `TriggerWebhooks` is deprecated in favor of it.

## [1.14.4] - 2026-08-07

### New Features
#### Send options
- **Exported `SendOptions` fields** — `SendOptions.SkipPhoneValidation` and `SendOptions.DeviceActiveWithin` are now exported and carry `query` struct tags, so callers can inspect or build options directly instead of only through option functions.
- **Added `SendOptions.Validate()`** — validates that `DeviceActiveWithin` is between 1 and `math.MaxInt32` when set; returns an error wrapping `ErrValidationFailed` when out of range. `WithDeviceActiveWithin` now clamps values above `math.MaxInt32` to `math.MaxInt32`.

### Bug Fixes
- **`Client.Send` rejects invalid send options before sending** — an out-of-range `DeviceActiveWithin` (for example `0`) now fails fast client-side with an error wrapping `ErrValidationFailed` instead of being sent to the gateway as an invalid request.