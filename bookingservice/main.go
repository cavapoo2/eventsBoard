package main

import (
	"flag"
	"time"

	"andy/booking_publish/bookingservice/listener"
	"andy/booking_publish/bookingservice/rest"
	"andy/booking_publish/lib/configuration"
	"andy/booking_publish/lib/msgqueue"
	msgqueue_amqp "andy/booking_publish/lib/msgqueue/amqp"
	"andy/booking_publish/lib/msgqueue/kafka"
	"andy/booking_publish/lib/persistence/dblayer"

	"log"

	"github.com/Shopify/sarama"
	"github.com/streadway/amqp"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
func init() {

	//log.SetOutput(ioutil.Discard)

}

func main() {
	var eventListener msgqueue.EventListener
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("conf", "./configuration/config.json", "flag to set the path to the configuration json file")
	flag.Parse()

	//extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath)

	switch config.MessageBrokerType {
	case "amqp":
		log.Println("BookingService: messageBrokerType=amqp")
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		if err != nil {
			time.Sleep(10 * time.Second)
			conn, err = amqp.Dial(config.AMQPMessageBroker)
		}

		panicIfErr(err)

		eventListener, err = msgqueue_amqp.NewAMQPEventListener(conn, "events", "booking")
		panicIfErr(err)

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		panicIfErr(err)
	case "kafka":
		log.Println("BookingService: messageBrokerType=kafka")
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		panicIfErr(err)

		eventListener, err = kafka.NewKafkaEventListener(conn, []int32{})
		panicIfErr(err)

		eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
		panicIfErr(err)
	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	processor := listener.EventProcessor{eventListener, dbhandler}
	go processor.ProcessEvents()

	rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)
}
