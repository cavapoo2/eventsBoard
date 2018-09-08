package contracts

import "github.com/cavapoo2/eventsBoard/lib/persistence"

//UserCreatedEvent is emitted whenever a new user is created
type UserCreatedEvent struct {
	ID       string                `json:"id"`
	First    string                `json:"first"`
	Last     string                `json:"last"`
	Email    string                `json:"email"`
	Password string                `json:"password"`
	Age      int                   `json:"age"`
	Bookings []persistence.Booking `json:"bookings"`
}

func (u *UserCreatedEvent) EventName() string {
	return "userCreated"
}
