import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Switch, Route,HashRouter } from 'react-router-dom';
import { Login } from './components/Login';
import { Register } from './components/Register';
import { EventBookingFormContainer } from './components/event_booking_form_container';
import { EventListContainer } from "./components/event_list_container";
import { BookingListContainer } from "./components/booking_list_container";
import {ErrorLogin} from './components/error'
import {AdminLogin} from './components/admin_login'
import {AdminEvent} from './components/admin_events'
import {AdminEventsCreatedContainer} from './components/admin_events_created_container'

class App extends React.Component<{}, {}> {
  render() {
    const eventList = (props) => <EventListContainer eventServiceURL="http://localhost:8181" {...props}/>;
    const userBookings = ({match}:any,props) => <BookingListContainer bookingServiceURL="http://localhost:8182" userid={match.params.userid} {...props}/>
    const eventBooking = ({ match }: any) => <EventBookingFormContainer userID={match.params.userid} eventID={match.params.id} eventServiceURL="http://localhost:8181"
      bookingServiceURL="http://localhost:8182" />;
    const adminEvent = ({match}:any,props) => <AdminEvent userID={match.params.userid} {...props}/>  
    const adminEventLogin = (props) => <AdminEvent {...props}/>  
    const adminEventsCreated = ({match}) => <AdminEventsCreatedContainer eventServiceURL="http://localhost:8181" userID={match.params.userid}/>  

    return <HashRouter >
    <Switch>
        <Route exact={true} path="/" component={Login} />
        <Route path="/register" component={Register} />
        <Route path="/list" component={eventList} />
        <Route path="/error" component={ErrorLogin}/>
        <Route path="/bookings/:id/:userid" component={eventBooking}/>
        <Route path="/userbookings/:userid" component={userBookings} />
        <Route exact={true} path="/admin" component={AdminLogin} />
        <Route exact={true} path="/admin/event" component={adminEventLogin}/>
        <Route exact={true} path="/admin/event/:userid" component={adminEvent}/>
        <Route path="/admin/eventsCreated/:userid" component={adminEventsCreated}/>

        </Switch>
    </HashRouter>
  }
}

ReactDOM.render(<App />, document.getElementById('root'));
