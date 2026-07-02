package smsgateway

type JWTScope = string

const (
	ScopeDevicesList   JWTScope = "devices:list"
	ScopeDevicesDelete JWTScope = "devices:delete"

	ScopeInboxList    JWTScope = "inbox:list"
	ScopeInboxRefresh JWTScope = "inbox:refresh"

	ScopeLogsRead JWTScope = "logs:read"

	ScopeMessagesCancel JWTScope = "messages:cancel"
	ScopeMessagesSend   JWTScope = "messages:send"
	ScopeMessagesRead   JWTScope = "messages:read"
	ScopeMessagesList   JWTScope = "messages:list"
	ScopeMessagesExport JWTScope = "messages:export"

	ScopeSettingsRead  JWTScope = "settings:read"
	ScopeSettingsWrite JWTScope = "settings:write"

	ScopeTokensManage JWTScope = "tokens:manage"

	ScopeWebhooksList   JWTScope = "webhooks:list"
	ScopeWebhooksWrite  JWTScope = "webhooks:write"
	ScopeWebhooksDelete JWTScope = "webhooks:delete"
)
