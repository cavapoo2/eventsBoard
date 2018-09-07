import * as React from "react";
import {Booking} from "../model/event";
import { Button, FormGroup, FormControl, ControlLabel, FormControlProps, HelpBlock, FormControlFeedback } from "react-bootstrap";
import { createBrowserHistory } from "history";

export interface AdminHallItemProps {
    index:number;
  //  handler:(e : React.FormEvent<FormControlProps>) => any
    //valid:(c:number) => boolean;
    cb:(e:React.FormEvent<FormControlProps>,i: number) => any;

}
export interface AdminHallItemState {
    Name:string;
    Location:string;
    Capacity:number;
}

export class AdminHallItem extends React.Component<AdminHallItemProps, AdminHallItemState> {
    constructor(props: AdminHallItemProps){
        super(props);
        this.handleChange = this.handleChange.bind(this);
        this.handleValid = this.handleValid.bind(this);
        
        this.state={
            Name:'',
            Location:'',
            Capacity:0,
        };
    }
    handleChange(e: React.FormEvent<FormControlProps>) {
	    let name = e.currentTarget.name;
        let val = e.currentTarget.value;
        this.setState({...this.state,[name]:val});
       // console.log('item=',this.state)
        

    }
    handleValid(v:number) {
        return v > 0 ? true : false;
    }
    render() {
        return (
            <div>
				<FormGroup controlId="nameHall" bsSize="small">
							<ControlLabel>Hall name</ControlLabel>
							<FormControl
								value={this.state.Name}
								onChange={e => {this.handleChange(e);this.props.cb(e,this.props.index)}}
								type="string"
								name="Name"
							/>
						</FormGroup>
						<FormGroup controlId="locHall" bsSize="small">
							<ControlLabel>Hall Location</ControlLabel>
							<FormControl
								value={this.state.Location}
								onChange={e =>{this.handleChange(e); this.props.cb(e,this.props.index)}}
								type="string"
								name="Location"
							/>
						</FormGroup>
						<FormGroup controlId="nameHall" bsSize="small">
                            <ControlLabel> Capacity </ControlLabel>
							<FormControl
								value={this.state.Capacity}
								onChange={e => {this.handleChange(e); this.props.cb(e,this.props.index)}}
								type="number"
								name="Capacity"
							/>
							<FormControl.Feedback />
							<HelpBlock> {this.handleValid(this.state.Capacity) ? "" : "must have greater than 0 seats"}</HelpBlock>

						</FormGroup>
        </div>
        )     
    }
}