import * as React from "react";
import {Link} from "react-router-dom";
import {Event} from "../model/event";

export interface EventListItemProps {
    userID:string;
    event: Event;
    selected?: boolean;

    onBooked: () => any;
}

export class EventListItem extends React.Component<EventListItemProps, {}> {
    render() {
        console.log('EventListItem userid=',this.props.userID)
        const start = new Date(this.props.event.StartDate * 1000);
        const end = new Date(this.props.event.EndDate * 1000);

        const locationName = this.props.event.Location ? this.props.event.Location.Name : "unknown";

        console.log(this.props.event);

        return <tr>
            <td>{this.props.event.Name}</td>
            <td>{locationName}</td>
            <td>{start.toLocaleDateString()}</td>
            <td>{end.toLocaleDateString()}</td>
            <td><Link to={`/bookings/${this.props.event.ID}/${this.props.userID}/bookings`} className="btn btn-primary">Book now!</Link></td>
        </tr>
    }
}