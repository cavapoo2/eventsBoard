import * as React from "react";
import { EventListItem } from "./event_list_item";
import { Event } from "../model/event";
import { Link } from 'react-router-dom';
export interface EventListProps {
    userID: string;
    name:string;
    events: Event[];
    onEventBooked: (e: Event) => any;
}

//<td><Link to={`/bookings/${this.props.event.ID}/${this.props.userID}/bookings`}
export class EventList extends React.Component<EventListProps, {}> {
    public render() {
       // console.log('EventList userid=', this.props.userID)
        const items = this.props.events.map(event =>
            <EventListItem key={event.ID} userID={this.props.userID} event={event} onBooked={() => this.props.onEventBooked(event)} />
        );

        return <div>
            <nav className="navbar navbar-expand-md navbar-dark bg-dark">
                <div className="navbar-text"> <strong>Available Events</strong></div>
                <div className="nav navbar-nav ml-auto">
                    <Link to={`/userbookings/${this.props.userID}`}>See bookings for {this.props.name} </Link>
                </div>
            </nav>
            <table className="table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Where</th>
                        <th colSpan={2}>When (start/end)</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {items}
                </tbody>
            </table>
        </div>
    }
}