INTRO
-----

This is my first react/typescript app still work in progress.
Its a simple events/ booking system. Users can sign up (register) and log in and see current events. Admin users create events for normal users to see. Normal users can then book events. 
The front end is using typescript and the backend is golang. The backend is using mongo db (2 instances) to persist as well as rabbitmq to send messages to either database. This is also Dockerized. 

INSTALL
-------
`go get -u github.com/cavapoo2/eventsBoard`

cd into tools and run:

`go build builder.go`

then cd .. and run 

`tools\builder`

that should build the main programs for golang back end the booking service and event service

cd into booking_front and run (note you need node installed, I'm using node v8.11.3 and npm 5.6.0 on ubuntu 16.04)

`npm install` 

to set up the back end (only tested locally)

first make sure you have rabbitmq setup (this is for ubuntu via docker, so yes you need docker installed)
`
docker run --detach \
--name rabbitmq \
-p 5672:5672 \
-p 15672:15672 \
rabbitmq:3-management`

make sure mongo is setup/installed.

in one terminal do :

`sudo mongod --dbpath /var/lib/mongodb --port 27017`

in another terminal do:

`mongod --dbpath ~/go/src/andy/booking/bookingservice/db --port 27018`

you now have 2 instances of mongod running . The first mongod will use /etc/mongo.conf settings for dbpath. you can optionally chose dbpath (with --dbpath) if you want just like the second command did. you can also use the `sudo service mongod start` approach.

in the eventService folder (in another terminal) do :

`./main` 

to start the backend eventservice. Also in the bookingservice folder (in another terminal) do:

`./main -conf=config.json`

to start the backend booking service.


this config.json sets up the bookingservice to use mongo port 27018
and the rest interface for the booking service is localhost:8182

open 2 more terminals one for (both for mongo cli)  

`mongo --port 27017`

and another

`mongo --port 27018`

you can now start the front end (cd into booking_front) :

`npm run start`

ALTERNATIVE DOCKER SETUP
-------------------------

make sure the binary exe's in both bookingservice and eventservice are name to these names rather than being called main.

the Dockerfile in each of the bookingservice and eventservice folders are as follows:

#event service docker file (this is the dockerfile)

`FROM debian:jessie
COPY eventservice /eventservice
RUN useradd eventservice
USER eventservice
ENV LISTEN_URL=0.0.0.0:8181
EXPOSE 8181
CMD ["/eventservice"]`

build the image like this in eventService folder like this:

`docker image build -t myevents/eventservice .`

likewise in the bookingservice folder (this is the dockerfile)

`FROM debian:jessie
COPY bookingservice /bookingservice
RUN useradd bookingservice
USER bookingservice
ENV LISTEN_URL=0.0.0.0:8181
EXPOSE 8181
CMD ["/bookingservice"]`

build like this (in the bookingservice folder)

`docker image build -t myevents/bookingservice .`

add a container network

`docker network create myevents`

run rabbitmq like this

`docker container run -d --name rabbitmq --network myevents
rabbitmq:3-management`

run events-db like this

`docker container run -d --name events-db --network myevents mongo`

bookings-db like this (make sure its a seperate instance, hence 20718)

`docker container run -d -p 20718:20718 --name bookings-db --network myevents mongo`

run eventservice like this

`docker container run \
--detach \
--name events \
--network myevents \
-e AMQP_BROKER_URL=amqp://guest:guest@rabbitmq:5672/ \
-e MONGO_URL=mongodb://events-db/events \
-p 8181:8181 \
myevents/eventservice`

run the bookingservice like this

`docker container run --detach --name bookings --network myevents -e AMQP_BROKER_URL=amqp://guest:guest@rabbitmq:5672/ -e MONGO_URL=mongodb://bookings-db/bookings -p 8182:8181 myevents/bookingservice`

run the booking_front end like this

`docker container build -t myevents/frontend .`

then

`docker container run --name frontend -p 80:80 myevents/frontend`	 

ALTERNATIVE WITH DOCKER COMPOSE
--------------------------------

Alternatively rather than using docker as above (whuch is good way to test its all working) can instead use the docker-compose which will build and run the whole lot in one go

`docker-compose up -d`

CREATING SOME ADMIN USERS
--------------------------



#use the simpleEvent.go in tools to build some json input files for testing the backend. 
#use the bash scripts to send GET/POST messages.
#use the builder to build project go files (see other directory (simple builder)

now use the bash script to create an even (newEvent). and another to create a new user (newUser). check the mongo cli
both these scripts talk to endpoint localhost:8181. these events are written to mongo db (27017) and via amqp its replicated
to the other mongo db at 27018. 
now send a makeBooking . this talks to end point localhost:8182, and updates the mongo db at 27018 only.


#check its all running with docker ps -a
#use the bash script 
newEvent outx.json 
#to send events
#also bash sript to send new user
newUser userx.json
#then you can make booking via
makeBooking bash script.

#you can view this by going to http://localhost



TODO (LOTS!)
------------

There is lots more work that could be done with this such as :
(1) Test more, add testing, Unit testing , improve code structures.
(2) improve user interface. Its using very simple Bootstrap right now.
(3) Think how this can be scaled better. 
(4) improve authentication. Its really basic email and password so far.
(5) add travis or CI
(6) improve logging
(7) maybe use redux. 
(8) monitoring (data dog)


    © 2018 GitHub, Inc.
    Terms
    Privacy
    Security
    Status
    Help

    Contact GitHub
    Pricing
    API
    Training
    Blog
    About

