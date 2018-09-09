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
```
docker run --detach \

--name rabbitmq \

-p 5672:5672 \

-p 15672:15672 \

rabbitmq:3-management

```

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

this is the event service docker file

```
FROM debian:jessie

COPY eventservice /eventservice

RUN useradd eventservice

USER eventservice

ENV LISTEN_URL=0.0.0.0:8181

EXPOSE 8181

CMD ["/eventservice"]
```

build the image like this in eventservice folder:

`docker image build -t myevents/eventservice .`

likewise in the bookingservice folder, the docker file is as follows:

```
FROM debian:jessie

COPY bookingservice /bookingservice

RUN useradd bookingservice

USER bookingservice

ENV LISTEN_URL=0.0.0.0:8181

EXPOSE 8181

CMD ["/bookingservice"]
```

build like this (in the bookingservice folder)

`docker image build -t myevents/bookingservice .`

add a container network.

`docker network create myevents`

run rabbitmq like this

`docker container run -d --name rabbitmq --network myevents
rabbitmq:3-management`

run events-db like this

`docker container run -d --name events-db --network myevents mongo`

bookings-db like this (make sure its a seperate instance, hence 20718)

`docker container run -d -p 20718:20718 --name bookings-db --network myevents mongo`

run eventservice like this
```
docker container run \
--detach \
--name events \
--network myevents \
-e AMQP_BROKER_URL=amqp://guest:guest@rabbitmq:5672/ \
-e MONGO_URL=mongodb://events-db/events \
-p 8181:8181 \
myevents/eventservice
```

run the bookingservice like this

`docker container run --detach --name bookings --network myevents -e AMQP_BROKER_URL=amqp://guest:guest@rabbitmq:5672/ -e MONGO_URL=mongodb://bookings-db/bookings -p 8182:8181 myevents/bookingservice`

run the booking_front end like this (cd into booking_front)

`docker container build -t myevents/frontend .`

then

`docker container run --name frontend -p 80:80 myevents/frontend`	 

you can go to:

localhost

in your web browser.

ALTERNATIVE WITH DOCKER COMPOSE
--------------------------------

Alternatively rather than using docker as above (whuch is good way to test its all working) can instead use the docker-compose which will build and run the whole lot in one go. Note good to do following if starting from fresh

`docker kill $(docker ps -q)`

`docker system prune -a`


`docker-compose up -d`

you can use this command to view the mongo database if using these docker containers

`docker exec -it container_id mongo`

this will bring you to the mongo shell where you can run 

`show dbs`

`use some_db`

`show collections`

`db.some_db.find().pretty()`



CREATING SOME ADMIN USERS
--------------------------

Admin user creation is not completed yet. I think typically owner of site would vet admin user and then the site owner would setup a password for the admin user. To get around this for now just use the admin bash script in scripts

`adminUser adminuser0.json`

then log in as an admin user using these credentials. To login via and admin, first go to the localhost:8080 or just localhost if using dockered version. The hash should appear at the end (localhost:8080/#) and just append admain like this
(localhost:8080/#/admin) . When you log in as an admin user you can then create events. Normal users can only book events.


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

(9) get a production build working for the font_end


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

