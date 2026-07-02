package smsgateway_test

import (
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func TestInboxRefreshRequest_ResolveWebhookDelivery(t *testing.T) {
	batch := smsgateway.WebhookDeliveryBatch
	individual := smsgateway.WebhookDeliveryIndividual
	disabled := smsgateway.WebhookDeliveryDisabled

	tests := []struct {
		name    string
		request smsgateway.InboxRefreshRequest
		want    smsgateway.WebhookDelivery
	}{
		{
			name:    "WebhookDelivery set to Batch",
			request: smsgateway.InboxRefreshRequest{WebhookDelivery: &batch},
			want:    smsgateway.WebhookDeliveryBatch,
		},
		{
			name:    "WebhookDelivery set to Individual",
			request: smsgateway.InboxRefreshRequest{WebhookDelivery: &individual},
			want:    smsgateway.WebhookDeliveryIndividual,
		},
		{
			name:    "WebhookDelivery set to Disabled overrides TriggerWebhooks",
			request: smsgateway.InboxRefreshRequest{WebhookDelivery: &disabled, TriggerWebhooks: true},
			want:    smsgateway.WebhookDeliveryDisabled,
		},
		{
			name:    "TriggerWebhooks true without WebhookDelivery",
			request: smsgateway.InboxRefreshRequest{TriggerWebhooks: true},
			want:    smsgateway.WebhookDeliveryIndividual,
		},
		{
			name:    "Default (neither set)",
			request: smsgateway.InboxRefreshRequest{},
			want:    smsgateway.WebhookDeliveryDisabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.request.ResolveWebhookDelivery(); got != tt.want {
				t.Errorf("ResolveWebhookDelivery() = %v, want %v", got, tt.want)
			}
		})
	}
}
