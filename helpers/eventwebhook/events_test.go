package eventwebhook

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseEvents(t *testing.T) {
	errWrongTypeMapping := fmt.Errorf("wrong type mapping")

	testCases := []struct {
		name          string
		payload       string
		expectedType  Event
		expectedError error
	}{
		{
			name:         "valid processed",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "pool": { "name": "new_MY_test", "id": 210 }, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"processed", "category":"cat facts", "sg_event_id":"rbtnWrG1DVDGGGFHFyun0A==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.000000000000000000000" } ]`,
			expectedType: &ProcessedEvent{},
		},
		{
			name:          "wrong expected type; processed vs dropped",
			payload:       `[ { "email":"example@test.com", "timestamp":1513299569, "pool": { "name": "new_MY_test", "id": 210 }, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"processed", "category":"cat facts", "sg_event_id":"rbtnWrG1DVDGGGFHFyun0A==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.000000000000000000000" } ]`,
			expectedType:  &DroppedEvent{},
			expectedError: errWrongTypeMapping,
		},
		{
			name:         "valid dropped",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"dropped", "category":"cat facts", "sg_event_id":"zmzJhfJgAfUSOW80yEbPyw==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0", "reason":"Bounced Address", "status":"5.0.0" } ]`,
			expectedType: &DroppedEvent{},
		},
		{
			name:         "valid delivered",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"delivered", "category":"cat facts", "sg_event_id":"rWVYmVk90MjZJ9iohOBa3w==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0", "response":"250 OK" } ]`,
			expectedType: &DelieveredEvent{},
		},
		{
			name:         "valid deferred",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"deferred", "category":"cat facts", "sg_event_id":"t7LEShmowp86DTdUW8M-GQ==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0", "response":"400 try again later", "attempt":"5" } ]`,
			expectedType: &DeferredEvent{},
		},
		{
			name:         "valid bounce",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"bounce", "category":"cat facts", "sg_event_id":"6g4ZI7SA-xmRDv57GoPIPw==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0", "reason":"500 unknown recipient", "status":"5.0.0", "type":"bounce" } ]`,
			expectedType: &BounceEvent{},
		},
		{
			name:         "valid open",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"open", "category":"cat facts", "sg_event_id":"FOTFFO0ecsBE-zxFXfs6WA==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0", "useragent":"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP; .NET CLR 1.1.4322; .NET CLR 2.0.50727)", "ip":"255.255.255.255" } ]`,
			expectedType: &OpenEvent{},
		},
		{
			name:         "valid click",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"click", "category":"cat facts", "sg_event_id":"kCAi1KttyQdEKHhdC-nuEA==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0", "useragent":"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP; .NET CLR 1.1.4322; .NET CLR 2.0.50727)", "ip":"255.255.255.255", "url":"http://www.sendgrid.com/" } ]`,
			expectedType: &ClickEvent{},
		},
		{
			name:         "valid spam report",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"spamreport", "category":"cat facts", "sg_event_id":"37nvH5QBz858KGVYCM4uOA==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0" } ]`,
			expectedType: &SpamReportEvent{},
		},
		{
			name:         "valid unsubscribe",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"unsubscribe", "category":"cat facts", "sg_event_id":"zz_BjPgU_5pS-J8vlfB1sg==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0" } ]`,
			expectedType: &UnsubscribeEvent{},
		},
		{
			name:         "valid group unsubscribe",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"group_unsubscribe", "category":"cat facts", "sg_event_id":"ahSCB7xYcXFb-hEaawsPRw==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0", "useragent":"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP; .NET CLR 1.1.4322; .NET CLR 2.0.50727)", "ip":"255.255.255.255", "url":"http://www.sendgrid.com/", "asm_group_id":10 } ]`,
			expectedType: &GroupUnsubscribeEvent{},
		},
		{
			name:         "valid group resubscribe",
			payload:      `[ { "email":"example@test.com", "timestamp":1513299569, "smtp-id":"<14c5d75ce93.dfd.64b469@ismtpd-555>", "event":"group_resubscribe", "category":["cat","facts"], "sg_event_id":"w_u0vJhLT-OFfprar5N93g==", "sg_message_id":"14c5d75ce93.dfd.64b469.filter0001.16648.5515E0B88.0", "useragent":"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP; .NET CLR 1.1.4322; .NET CLR 2.0.50727)", "ip":"255.255.255.255", "url":"http://www.sendgrid.com/", "asm_group_id":10 } ]`,
			expectedType: &GroupResubscribeEvent{},
		},
	}

	for _, tc := range testCases {
		evts := []json.RawMessage{}
		if err := json.Unmarshal([]byte(tc.payload), &evts); err != nil {
			t.Error(err)
		}

		for _, e := range evts {
			ge := &GenericEvent{}
			if err := json.Unmarshal(e, ge); err != nil {
				t.Error(err)
			}

			specificType, err := TypeByEventName(ge.EventType, e)
			if err != nil {
				t.Error(err)
			}

			// checking the type is what it should be
			a := reflect.TypeOf(specificType).Elem()
			b := reflect.TypeOf(tc.expectedType).Elem()
			if a != b {
				if tc.expectedError != nil && tc.expectedError == errWrongTypeMapping {
					continue
				}

				t.Errorf("types do not match: %s != %s", a.Name(), b.Name())
			}

			// checking the values match up
			m := make(map[string]interface{})
			if err = json.Unmarshal(e, &m); err != nil {
				t.Error(err)
			}

			stValue := reflect.Indirect(reflect.ValueOf(specificType))
			stType := a
			for i := 0; i < a.NumField(); i++ {
				field := stValue.Field(i)
				fieldValue := field.Interface()

				fieldType := stType.Field(i)
				fieldJSONName := strings.Split(fieldType.Tag.Get("json"), ",")[0]

				// skipping embeded struct for now
				if !fieldType.Anonymous && fmt.Sprint(m[fieldJSONName]) != fmt.Sprint(fieldValue) {
					fmt.Println(m)
					t.Errorf("%s: %v != %v", fieldJSONName, m[fieldJSONName], fieldValue)
				}
			}

		}

	}
}
