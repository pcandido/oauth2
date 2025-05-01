package store

type Event struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end"`
	UserId string
}

var events []Event = []Event{
	{
		ID:     "1",
		Title:  "Meeting with client",
		Start:  "2025-05-01T10:00:00Z",
		End:    "2025-05-01T12:00:00Z",
		UserId: "user1",
	},
	{
		ID:     "2",
		Title:  "Product review",
		Start:  "2025-05-01T13:00:00Z",
		End:    "2025-05-01T15:00:00Z",
		UserId: "user1",
	},
	{
		ID:     "3",
		Title:  "Retrospective",
		Start:  "2025-05-01T16:00:00Z",
		End:    "2025-05-01T16:00:00Z",
		UserId: "user1",
	},
	{
		ID:     "4",
		Title:  "Team lunch",
		Start:  "2025-05-01T12:00:00Z",
		End:    "2025-05-01T13:00:00Z",
		UserId: "user2",
	},
	{
		ID:     "5",
		Title:  "Product review",
		Start:  "2025-05-01T13:00:00Z",
		End:    "2025-05-01T15:00:00Z",
		UserId: "user2",
	},
}

func GetEventsByUserId(userId string) []Event {
	var userEvents []Event
	for _, event := range events {
		if event.UserId == userId {
			userEvents = append(userEvents, event)
		}
	}
	return userEvents
}
