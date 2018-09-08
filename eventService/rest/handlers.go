package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cavapoo2/eventsBoard/contracts"
	"github.com/cavapoo2/eventsBoard/lib/msgqueue"
	"github.com/cavapoo2/eventsBoard/lib/persistence"

	"github.com/gorilla/mux"
)

type eventServiceHandler struct {
	dbhandler    persistence.DatabaseHandler
	eventEmitter msgqueue.EventEmitter
}

func newEventHandler(databasehandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) *eventServiceHandler {
	return &eventServiceHandler{
		dbhandler:    databasehandler,
		eventEmitter: eventEmitter,
	}
}

/*
func (eh *eventServiceHandler) checkAdminUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pass, ok := vars["password"] //TODO should be encrypting passwords!
	if !ok {
		w.WriteHeader(404)
		fmt.Fprintf(w, "No password found")
		return
	}
	//just use the server to read a file for now. //TODO make more secure

}*/

//note in practive the password would be encrypted on client side, then decrypted here and send to db
func (eh *eventServiceHandler) findUserEmailPassHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		fmt.Fprint(w, `No firstname found`)
		return
	}
	pass, ok1 := vars["password"]
	if !ok1 {
		fmt.Fprint(w, `No secondname found`)
		return

	}

	u := persistence.User{}
	u, err := eh.dbhandler.FindUserEmailPass(email, pass)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Error occured %s", err)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&u)

}

//ditto as above with regsrds to security
func (eh *eventServiceHandler) verifyAdminUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		fmt.Fprint(w, `No firstname found`)
		return
	}
	pass, ok1 := vars["password"]
	if !ok1 {
		fmt.Fprint(w, `No secondname found`)
		return

	}

	u := persistence.AdminUser{}
	u, err := eh.dbhandler.FindAdminUserEmailPass(email, pass)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Error occured %s", err)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&u)

}

func (eh *eventServiceHandler) findUserHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	first, ok := vars["firstname"]
	if !ok {
		fmt.Fprint(w, `No firstname found`)
		return
	}
	second, ok1 := vars["secondname"]
	if !ok1 {
		fmt.Fprint(w, `No secondname found`)
		return

	}

	u := persistence.User{}
	u, err := eh.dbhandler.FindUser(first, second)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Error occured %s", err)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&u)

}

func (eh *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		fmt.Fprint(w, `No search criteria found, you can either search by id via /id/4
						to search by name via /name/coldplayconcert`)
		return
	}

	searchkey, ok := vars["search"]
	if !ok {
		fmt.Fprint(w, `No search keys found, you can either search by id via /id/4
						to search by name via /name/coldplayconcert`)
		return
	}

	var event persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchkey)
	case "id":
		id, err := hex.DecodeString(searchkey)
		if nil == err {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Error occured %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying to find all available events %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying encode events to JSON %s", err)
	}
}

func (eh *eventServiceHandler) oneEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, ok := vars["eventID"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, "missing route parameter 'eventID'")
		return
	}

	eventIDBytes, _ := hex.DecodeString(eventID)
	event, err := eh.dbhandler.FindEvent(eventIDBytes)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "event with id %s was not found", eventID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding event data %s", err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting event %s", err)
		return
	}

	msg := contracts.EventCreatedEvent{
		ID:         hex.EncodeToString(id),
		Name:       event.Name,
		Start:      time.Unix(event.StartDate, 0),
		End:        time.Unix(event.EndDate, 0),
		LocationID: string(event.Location.ID),
	}
	eh.eventEmitter.Emit(&msg)

	w.Header().Set("Content-Type", "application/json;charset=utf8")

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) allLocationsHandler(w http.ResponseWriter, r *http.Request) {
	locations, err := eh.dbhandler.FindAllLocations()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "could not load locations: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")

	err = json.NewEncoder(w).Encode(locations)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying encode events to JSON %s", err)
	}

}

func (eh *eventServiceHandler) newLocationHandler(w http.ResponseWriter, r *http.Request) {
	location := persistence.Location{}
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "request body could not be unserialized to location: %s", err)
		return
	}

	persistedLocation, err := eh.dbhandler.AddLocation(location)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "could not persist location: %s", err)
	}

	msg := contracts.LocationCreatedEvent{
		ID:      string(persistedLocation.ID),
		Name:    persistedLocation.Name,
		Address: persistedLocation.Address,
		Country: persistedLocation.Country,
		Halls:   persistedLocation.Halls,
	}
	eh.eventEmitter.Emit(&msg)

	w.Header().Set("Content-Type", "application/json;charset=utf8")

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&persistedLocation)
}

//add this not sure how it fits in yet
func (eh *eventServiceHandler) newUserHandler(w http.ResponseWriter, r *http.Request) {
	//in post pass in id, first and last name,age and bookings

	user := persistence.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "request body could not be unserialized to user: %s", err)
		return
	}

	persistedUser, err := eh.dbhandler.AddUser(user)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "could not persist user: %s", err)
	}

	msg := contracts.UserCreatedEvent{
		ID: hex.EncodeToString(persistedUser),
		//ID:       string(user.ID),
		First:    user.First,
		Last:     user.Last,
		Email:    user.Email,
		Password: user.Password, //clearly this is not something that is passed around and probably hashed
		Age:      user.Age,
		Bookings: user.Bookings,
	}
	err = eh.eventEmitter.Emit(&msg)
	if err != nil {
		log.Printf("Issue Emitting +%v", err)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&persistedUser)
}

//add this not sure how it fits in yet
func (eh *eventServiceHandler) newAdminUserHandler(w http.ResponseWriter, r *http.Request) {
	//in post pass in id, first and last name,age and bookings

	user := persistence.AdminUser{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "request body could not be unserialized to admin user: %s", err)
		return
	}

	persistedUser, err := eh.dbhandler.AddAdminUser(user)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "could not persist admin user: %s", err)
	}

	msg := contracts.AdminUserCreatedEvent{
		ID: hex.EncodeToString(persistedUser),
		//ID:       string(user.ID),
		First:    user.First,
		Last:     user.Last,
		Email:    user.Email,
		Password: user.Password, //clearly this is not something that is passed around and probably hashed
		Company:  user.Company,
		Events:   user.Events,
	}
	err = eh.eventEmitter.Emit(&msg)
	if err != nil {
		log.Printf("Issue Emitting +%v", err)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&persistedUser)
}

func (eh *eventServiceHandler) addEventForAdminUser(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	userID, ok := routeVars["userID"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, "missing route parameter 'userID")
		return
	}
	fmt.Println(userID)
	userIDMongo, _ := hex.DecodeString(userID)
	/*
		adminUser, err := eh.dbhandler.FindAdminUser(userIDMongo)
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "bo for %s could not be loaded: %s", userIDMongo, err)
			return
		}*/

	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding admin event data %s", err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting event to public events database %s", err)
		return
	}
	//emit this to the other db
	msg := contracts.EventCreatedEvent{
		ID:         hex.EncodeToString(id),
		Name:       event.Name,
		Start:      time.Unix(event.StartDate, 0),
		End:        time.Unix(event.EndDate, 0),
		LocationID: string(event.Location.ID),
	}
	eh.eventEmitter.Emit(&msg)

	err = eh.dbhandler.AddEventForAdminUser(userIDMongo, event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting event for admin %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	err = json.NewEncoder(w).Encode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying to encode events to JSON %s", err)
	}

}

func (eh *eventServiceHandler) eventsForAdminUser(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	userID, ok := routeVars["userID"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, "missing route parameter 'userID")
		return
	}
	userIDMongo, _ := hex.DecodeString(userID)
	events, err := eh.dbhandler.FindEventsAdminUser(userIDMongo)

	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "events for %s admin user could not be loaded: %s", userIDMongo, err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying to encode events to JSON %s", err)
	}

}
