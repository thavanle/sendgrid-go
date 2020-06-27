package eventwebhook

type Events []Event

type Event interface {
}

type CommonEvent struct {
	Email     string `json:"email"`
	Timestamp int64  `json:"timestamp"`
	Event     string `json:"event"`
}

// ProcessedEvent ...
type ProcessedEvent struct {
	CommonEvent
	SMTPID                string
	SGEventID             string                 `json:"sg_event_id"`
	SGMessageID           string                 `json:"sg_message_id"`
	Category              interface{}            `json:"category"` //string, array[string]
	ASMGroupID            int                    `json:"asm_group_id"`
	UniqueArgs            map[string]interface{} `json:"unique_args"`
	MarketingCampaignID   int                    `json:"marketing_campaign_id"`
	MarketingCampaignName string                 `json:"marketing_campaign_name"`
	Pool                  map[string]interface{} `json:"pool"`
}

// DroppedEvent ...
type DroppedEvent struct {
	CommonEvent
}

// DelieveredEvent ...
type DelieveredEvent struct {
	CommonEvent
}

//  DeferredEvent ...
type DeferredEvent struct {
	CommonEvent
}

// BounceEvent ...
type BounceEvent struct {
	CommonEvent
}

// BlockedEvent ...
type BlockedEvent struct {
	CommonEvent
}

// OpenEvent ...
type OpenEvent struct {
	CommonEvent
}

// ClickEvent ...
type ClickEvent struct {
	CommonEvent
}

// SpamReportEvent ...
type SpamReportEvent struct {
	CommonEvent
}

// UnsubscribeEvent ...
type UnsubscribeEvent struct {
	CommonEvent
}

// GroupUnsubscribeEvent ...
type GroupUnsubscribeEvent struct {
	CommonEvent
}

// GroupResubscribeEvent ...
type GroupResubscribeEvent struct {
	CommonEvent
}
