import * as React from "react";
import {Event} from "../model/event";

export interface EventListItemProps {
    userID:string;
    event: Event;
    selected?: boolean;

//    onBooked: () => any;
}

export class AdminEventListItem extends React.Component<EventListItemProps, {}> {
    render() {
        console.log('EventListItem userid=',this.props.userID)
        const start = new Date(this.props.event.StartDate * 1000);
        const end = new Date(this.props.event.EndDate * 1000);
        //const startTime = this.props.event.Location.

        const locationName = this.props.event.Location ? this.props.event.Location.Name : "unknown";

        console.log(this.props.event);

        return <tr>
            <td>{this.props.event.Name}</td>
            <td>{locationName}</td>
            <td>{start.toLocaleDateString()}</td>
            <td>{end.toLocaleDateString()}</td>
            <td>{this.props.event.Location.OpenTime}</td>
            <td>{this.props.event.Location.CloseTime}</td>
            <td>{this.props.event.Location.Address}</td>


        </tr>
    }
}