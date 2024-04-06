package requests

import "time"

type FeedingSessionPerDayRequest struct {
	Day                  time.Weekday `json:"day"`
	FeedingSessionNumber string       `json:"feeding_session_number"`
}
