package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

//note had to combine commands to get them to work (proably is better way!)
//		"cd /home/febe/go/src/andy/builder/test && go build main.go",
func main() {
	installPaths := []string{

		"cd /home/febe/go/src/andy/booking_publish/lib/configuration && go install",
		"cd /home/febe/go/src/andy/booking_publish/lib/helper/amqp && go install",

		"cd /home/febe/go/src/andy/booking_publish/lib/helper/kafka && go install",
		"cd /home/febe/go/src/andy/booking_publish/lib/msgqueue/amqp && go install",
		"cd /home/febe/go/src/andy/booking_publish/lib/msgqueue/builder && go install",
		"cd /home/febe/go/src/andy/booking_publish/lib/msgqueue/kafka && go install",
		"cd /home/febe/go/src/andy/booking_publish/lib/msgqueue && go install",

		"cd /home/febe/go/src/andy/booking_publish/lib/persistence/mongolayer && go install",
		"cd /home/febe/go/src/andy/booking_publish/lib/persistence/dblayer && go install",
		"cd /home/febe/go/src/andy/booking_publish/lib/persistence && go install",

		"cd /home/febe/go/src/andy/booking_publish/bookingservice/listener && go install",
		"cd /home/febe/go/src/andy/booking_publish/bookingservice/rest && go install",
		"cd /home/febe/go/src/andy/booking_publish/contracts && go install",
		"cd /home/febe/go/src/andy/booking_publish/eventService/rest && go install",
		"cd /home/febe/go/src/andy/booking_publish/eventService/listner && go install",
	}

	buildPaths := []string{
		"cd /home/febe/go/src/andy/booking_publish/bookingservice && go build main.go",
		"cd /home/febe/go/src/andy/booking_publish/eventService && go build main.go",
	}
	for _, p := range installPaths {
		execCommand(p)
	}
	for _, p := range buildPaths {
		execCommand(p)
	}

}

func execCommand(c string) {
	//	out, err := exec.Command("/bin/bash", "-c", p).Output()
	cmd := exec.Command("/bin/bash", "-c", c)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ":  " + stderr.String())
	}

}
