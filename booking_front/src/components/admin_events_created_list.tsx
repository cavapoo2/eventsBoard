import * as React from "react";
import { AdminEventListItem } from "./admin_events_created_list_item";
import { Event } from "../model/event";
import { Link } from 'react-router-dom';
export interface AdminEventsCreatedListProps {
    userID: string;
   // name:string;
    events: Event[];
//    onEventBooked: (e: Event) => any;
}

//<td><Link to={`/bookings/${this.props.event.ID}/${this.props.userID}/bookings`}
export class AdminEventsCreatedList extends React.Component<AdminEventsCreatedListProps, {}> {
    public render() {
       // console.log('EventList userid=', this.props.userID)
       let i=0;
        const items = this.props.events.map(event =>
            <AdminEventListItem key={i++} userID={this.props.userID} event={event}  />
        );

        return <div>
            <nav className="navbar navbar-expand-md navbar-dark bg-dark">
                <div className="navbar-text"> <strong>Events Created</strong></div>
                <div className="nav navbar-nav ml-auto">
                    <Link to={`/admin/event/${this.props.userID}`}>Add more events</Link>
                </div>
                
            </nav>
            <table className="table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Where</th>
                        <th>Start Date</th>
                        <th>End Date</th>
                        <th>Start Time</th>
                        <th>End Time</th>
                        <th>Address</th>

                    </tr>
                </thead>
                <tbody>
                    {items}
                </tbody>
            </table>
        </div>
    }
}