import * as React from "react";
import {AdminEventsCreatedList} from "./admin_events_created_list";
import {Event} from "../model/event";
//import {Link} from 'react-router-dom';

export interface AdminEventsCreatedContainerProps {
    userID:string;
    eventServiceURL: string;
}

export interface AdminEventsCreatedContainerState {
    state: "loading"|"ready"|"saving"|"done"|"error";
    events: Event[];
}

export class AdminEventsCreatedContainer extends React.Component<AdminEventsCreatedContainerProps, AdminEventsCreatedContainerState> {
    constructor(p: AdminEventsCreatedContainerProps) {
        super(p);

	//this.handleSubmit = this.handleSubmit.bind(this);
        this.state = {
            state: "loading",
            events:[],
        };

        console.log('AdminEventsCreatedContainer url=',this.props.eventServiceURL);

        fetch(p.eventServiceURL + "/admin/eventsCreated/" + p.userID)
            .then<Event[]>(resp => resp.json())
            .then(events => {
                this.setState({
                    state: "ready",
                    events: events
                });
            })
    }

    render() {
        if (this.state.state === "loading") {
            console.log("state:loading")
            return <div>Loading...</div>;
        }
        console.log(this.state.events);
/*
        if (!this.state.event) {
            console.log('state:unknown')
            return <div>Unknown error</div>;
        }


        if (this.state.state === "done") {
            return <div>
            <nav className="navbar navbar-expand-md navbar-dark bg-dark">
                <div className="navbar-text"> <strong>Booking successfully completed</strong></div>
                <div className="nav navbar-nav ml-auto">
                    <Link to="/list">Back To Events</Link>
                 </div>
            </nav>
            </div>

//            return <div className="alert alert-success">Booking successfully completed!</div>
        }
*/
        return <AdminEventsCreatedList events={this.state.events} userID={this.props.userID} />
    }
/*
    handleSubmit(seats: number) {
        const url = this.props.bookingServiceURL + "/bookings/" + this.props.eventID + "/" + this.props.userID ;
       // const url = this.props.bookingServiceURL + "/" + this.props.eventID + "/" + this.props.userID ;
        const payload = {Seats: seats,Name:this.state.event.Name};
        console.log("state :saving")
        this.setState({
            event: this.state.event,
            state: "saving"
        });
        console.log('EventBookingFormContainer url=', url);
        console.log('payload=',JSON.stringify(payload));
        fetch(url, {method: "POST", body: JSON.stringify(payload)})
            .then(response => {
                this.setState({
                    event: this.state.event,
                    state: response.ok ? "done" : "error"
                });
            })
    }*/
}