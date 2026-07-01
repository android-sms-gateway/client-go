# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.14.4] - 2026-08-07

### New Features
#### Send options
- **Exported `SendOptions` fields** — `SendOptions.SkipPhoneValidation` and `SendOptions.DeviceActiveWithin` are now exported and carry `query` struct tags, so callers can inspect or build options directly instead of only through option functions.
- **Added `SendOptions.Validate()`** — validates that `DeviceActiveWithin` is between 1 and `math.MaxInt32` when set; returns an error wrapping `ErrValidationFailed` when out of range. `WithDeviceActiveWithin` now clamps values above `math.MaxInt32` to `math.MaxInt32`.

### Bug Fixes
- **`Client.Send` rejects invalid send options before sending** — an out-of-range `DeviceActiveWithin` (for example `0`) now fails fast client-side with an error wrapping `ErrValidationFailed` instead of being sent to the gateway as an invalid request.