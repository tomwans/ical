package ical

import (
	"bytes"
	"testing"
)

const calendarContents = `
BEGIN:VCALENDAR
METHOD:PUBLISH
VERSION:2.0
X-WR-CALNAME:TestCalendar
PRODID:-//Apple Inc.//Mac OS X 10.11.6//EN
X-APPLE-CALENDAR-COLOR:#CC73E1
X-WR-TIMEZONE:America/New_York
CALSCALE:GREGORIAN
BEGIN:VTIMEZONE
TZID:America/Los_Angeles
BEGIN:DAYLIGHT
TZOFFSETFROM:-0800
RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=2SU
DTSTART:20070311T020000
TZNAME:PDT
TZOFFSETTO:-0700
END:DAYLIGHT
BEGIN:STANDARD
TZOFFSETFROM:-0700
RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU
DTSTART:20071104T020000
TZNAME:PST
TZOFFSETTO:-0800
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
TRANSP:OPAQUE
DTEND;TZID=America/Los_Angeles:20160919T123000
ORGANIZER;CN="Thomas W":/aMTg0OTQ0MDYyMTg0OTQ0MPsrFnfsGgaWZGliH
 5cLau1T-aypWFNemxhPRSExqvHU/principal/
UID:08A7761E-78B7-422B-8FF2-34D68C9C53B4
DTSTAMP:20160919T030918Z
LOCATION:TestLocation\;KEY=VALUE
DESCRIPTION:TestNote\, TestNote2
SEQUENCE:0
SUMMARY:TestEvent
DTSTART;TZID=America/Los_Angeles:20160919T113000
X-APPLE-TRAVEL-ADVISORY-BEHAVIOR:AUTOMATIC
CREATED:20160919T030609Z
ATTENDEE;CN="Thomas W";CUTYPE=INDIVIDUAL;EMAIL="tw@ic
 loud.com";PARTSTAT=ACCEPTED;ROLE=CHAIR;X-CALENDARSERVER-DTSTAMP=20160919
 T030922Z:/aMTg0OTQ0MDYyMTg0OTQ0MPsrFnfsGgaWZGliH5cLau1T-aypWFNemxhPRSExq
 vHU/principal/
ATTENDEE;CN="tw@gmail.com";CUTYPE=INDIVIDUAL;EMAIL="tw@gmail.c
 om";PARTSTAT=NEEDS-ACTION;SCHEDULE-STATUS="1.1":mailto:tw@gmail.com
BEGIN:VALARM
X-WR-ALARMUID:D5FDDB21-9C2E-4353-81DA-963B79A19305
UID:D5FDDB21-9C2E-4353-81DA-963B79A19305
TRIGGER:-PT30M
X-APPLE-DEFAULT-ALARM:TRUE
ATTACH;VALUE=URI:Basso
ACTION:AUDIO
END:VALARM
END:VEVENT
END:VCALENDAR
`

func TestParseCalendar(t *testing.T) {
	buf := bytes.NewBufferString(calendarContents)
	dec := NewDecoder(buf)
	cal := &Token{}
	if err := dec.Decode(cal); err != nil {
		t.Error(err)
	}

	expected := 9
	got := len(cal.Subtokens)
	if got != expected {
		t.Errorf("got: %d subtokens, expected: %d", got, expected)
	}
}