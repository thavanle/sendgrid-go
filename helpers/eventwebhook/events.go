package eventwebhook

import (
	"encoding/json"
	"fmt"
)

type Events []Event

func (e *Events) UnmarshalJSON(data []byte) error {
	raw := []json.RawMessage{}
	err := json.Unmarshal(data, raw)
	if err != nil {
		return err
	}

	for _, rm := range raw {
		var b []byte
		err = rm.UnmarshalJSON(b)
		if err != nil {
			return err
		}

	}

	return nil
}

// Event ...
type Event interface {
	GetEventType() string
}

// GenericEvent ...
type GenericEvent struct {
	Email                 string                 `json:"email"`
	Timestamp             int64                  `json:"timestamp"`
	EventType             string                 `json:"event"`
	MarketingCampaignID   int                    `json:"marketing_campaign_id"`
	MarketingCampaignName string                 `json:"marketing_campaign_name"`
	UniqueArgs            map[string]interface{} `json:"unique_args"`
}

// GetEventType ...
func (ge GenericEvent) GetEventType() string {
	return ge.EventType
}

// TypeByEventName ...
func TypeByEventName(eventName string, rawMsg json.RawMessage) (Event, error) {
	switch eventName {
	case "processed":
		pe := &ProcessedEvent{}
		err := json.Unmarshal(rawMsg, pe)
		return pe, err
	case "dropped":
		de := &DroppedEvent{}
		err := json.Unmarshal(rawMsg, de)
		return de, err
	case "delivered":
		de := &DelieveredEvent{}
		err := json.Unmarshal(rawMsg, de)
		return de, err
	case "deferred":
		de := &DeferredEvent{}
		err := json.Unmarshal(rawMsg, de)
		return de, err
	case "bounce":
		be := &BounceEvent{}
		err := json.Unmarshal(rawMsg, be)
		return be, err
	case "open":
		oe := &OpenEvent{}
		err := json.Unmarshal(rawMsg, oe)
		return oe, err
	case "click":
		ce := &ClickEvent{}
		err := json.Unmarshal(rawMsg, ce)
		return ce, err
	case "spamreport":
		se := &SpamReportEvent{}
		err := json.Unmarshal(rawMsg, se)
		return se, err
	case "unsubscribe":
		ue := &UnsubscribeEvent{}
		err := json.Unmarshal(rawMsg, ue)
		return ue, err
	case "group_unsubscribe":
		ge := &GroupUnsubscribeEvent{}
		err := json.Unmarshal(rawMsg, ge)
		return ge, err
	case "group_resubscribe":
		ge := &GroupResubscribeEvent{}
		err := json.Unmarshal(rawMsg, ge)
		return ge, err
	}

	return nil, fmt.Errorf("unknown event: %s", eventName)
}

// EventCategory ...
type EventCategory []string

// UnmarshalJSON ...
func (ec *EventCategory) UnmarshalJSON(data []byte) error {
	arr := []string{}
	if err := json.Unmarshal(data, &arr); err != nil {
		*ec = []string{string(data)}
	}

	*ec = arr

	return nil
}

// ProcessedEvent ...
type ProcessedEvent struct {
	GenericEvent
	SMTPID      string        `json:"smtp-id"`
	SGEventID   string        `json:"sg_event_id"`
	SGMessageID string        `json:"sg_message_id"`
	Category    EventCategory `json:"category"` //string, array[string]
	ASMGroupID  int           `json:"asm_group_id"`
	//Pool        map[string]interface{} `json:"pool"`
}

// DroppedEvent ...
type DroppedEvent struct {
	GenericEvent
	SMTPID      string        `json:"smtp-id"`
	SGEventID   string        `json:"sg_event_id"`
	SGMessageID string        `json:"sg_message_id"`
	Reason      string        `json:"reason"`
	Category    EventCategory `json:"category"` //string, array[string]
	ASMGroupID  int           `json:"asm_group_id"`
}

// DelieveredEvent ...
type DelieveredEvent struct {
	GenericEvent
	SMTPID      string        `json:"smtp-id"`
	IP          string        `json:"ip"`
	SGEventID   string        `json:"sg_event_id"`
	SGMessageID string        `json:"sg_message_id"`
	Response    string        `json:"response"`
	TLS         bool          `json:"tls"`
	Category    EventCategory `json:"category"` //string, array[string]
	ASMGroupID  int           `json:"asm_group_id"`
}

// DeferredEvent ...
type DeferredEvent struct {
	GenericEvent
	SMTPID      string        `json:"smtp-id"`
	SGEventID   string        `json:"sg_event_id"`
	SGMessageID string        `json:"sg_message_id"`
	Reason      string        `json:"reason"`
	Category    EventCategory `json:"category"` //string, array[string]
	ASMGroupID  int           `json:"asm_group_id"`
	Attempt     int           `json:"attempt,string"`
}

// BounceEvent ...
type BounceEvent struct {
	GenericEvent
	SMTPID      string        `json:"smtp-id"`
	SGEventID   string        `json:"sg_event_id"`
	SGMessageID string        `json:"sg_message_id"`
	IP          string        `json:"ip"`
	Reason      string        `json:"reason"`
	Status      string        `json:"status"`
	TLS         bool          `json:"tls"`
	Category    EventCategory `json:"category"` //string, array[string]
	ASMGroupID  int           `json:"asm_group_id"`
}

// OpenEvent ...
type OpenEvent struct {
	GenericEvent
	UserAgent   string        `json:"useragent"`
	IP          string        `json:"ip"`
	SGEventID   string        `json:"sg_event_id"`
	SGMessageID string        `json:"sg_message_id"`
	Category    EventCategory `json:"category"` //string, array[string]
	ASMGroupID  int           `json:"asm_group_id"`
}

// ClickEvent ...
type ClickEvent struct {
	GenericEvent
	UserAgent   string        `json:"useragent"`
	IP          string        `json:"ip"`
	SGEventID   string        `json:"sg_event_id"`
	SGMessageID string        `json:"sg_message_id"`
	URL         string        `json:"url"`
	Category    EventCategory `json:"category"` //string, array[string]
	ASMGroupID  int           `json:"asm_group_id"`
}

// SpamReportEvent ...
type SpamReportEvent struct {
	GenericEvent
	SGEventID   string        `json:"sg_event_id"`
	SGMessageID string        `json:"sg_message_id"`
	Category    EventCategory `json:"category"` //string, array[string]

}

// UnsubscribeEvent ...
type UnsubscribeEvent struct {
	GenericEvent
	Category EventCategory `json:"category"` //string, array[string]
}

// GroupUnsubscribeEvent ...
type GroupUnsubscribeEvent struct {
	GenericEvent
	UserAgent   string `json:"useragent"`
	IP          string `json:"ip"`
	SGEventID   string `json:"sg_event_id"`
	SGMessageID string `json:"sg_message_id"`
	ASMGroupID  int    `json:"asm_group_id"`
}

// GroupResubscribeEvent ...
type GroupResubscribeEvent struct {
	GenericEvent
	UserAgent   string `json:"useragent"`
	IP          string `json:"ip"`
	SGEventID   string `json:"sg_event_id"`
	SGMessageID string `json:"sg_message_id"`
	ASMGroupID  int    `json:"asm_group_id"`
}
