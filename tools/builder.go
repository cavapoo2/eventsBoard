package main

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
)

//note had to combine commands to get them to work (proably is better way!)
//		"cd /home/febe/go/src/andy/builder/test && go build main.go",
func main() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	base := "cd " + gopath + "/src/github.com/cavapoo2/eventsBoard/"
	fmt.Println("base is ", base)
	installPaths := []string{

		base + "lib/configuration && go install",
		base + "lib/helper/amqp && go install",

		base + "lib/helper/kafka && go install",
		base + "lib/msgqueue/amqp && go install",
		base + "lib/msgqueue/builder && go install",
		base + "lib/msgqueue/kafka && go install",
		base + "lib/msgqueue && go install",

		base + "lib/persistence/mongolayer && go install",
		base + "lib/persistence/dblayer && go install",
		base + "lib/persistence && go install",

		base + "bookingservice/listener && go install",
		base + "bookingservice/rest && go install",
		base + "contracts && go install",
		base + "eventService/rest && go install",
		base + "eventService/listner && go install",
	}

	buildPaths := []string{
		base + "bookingservice && go build main.go",
		base + "eventService && go build main.go",
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
