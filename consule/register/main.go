package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gogocoding/micro-go-course/register/register/discovery"
	"github.com/gogocoding/micro-go-course/register/register/endpoint"
	"github.com/gogocoding/micro-go-course/register/register/service"
	"github.com/gogocoding/micro-go-course/register/register/transport"
	"github.com/google/uuid"
)

func main() {

	consulAddr := flag.String("consul.addr", "localhost", "consul address")
	consulPort := flag.Int("consul.port", 8500, "consul port")
	serviceName := flag.String("service.name", "register", "service name")
	serviceAddr := flag.String("service.addr", "localhost", "service addr")
	servicePort := flag.Int("service.port", 12312, "service port")

	flag.Parse()

	instanceId := *serviceName + "-" + uuid.New().String()
	client, err := discovery.NewDiscoveryClient(*consulAddr, *consulPort, discovery.NewAgentServiceRegistration(
		*serviceName, instanceId, "/health", *serviceAddr, *servicePort, nil))

	if err != nil {
		log.Printf("register service err : %s", err)
		os.Exit(-1)
	}

	errChan := make(chan error)

	srv := service.NewRegisterServiceImpl(client)

	endpoints := endpoint.RegisteryEndpoints{
		DiscoveryEndpoint:   endpoint.MakeDiscoveryEndpoint(srv),
		HealthCheckEndpoint: endpoint.MakeHealthCheckEndpoint(srv),
	}

	handler := transport.MakeHttpHandler(context.Background(), &endpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	client.Register(context.Background())
	err = <-errChan
	log.Printf("listen err: %s", err)
	client.Deregister(context.Background())

}

func init() {
	file := "./" + "register.log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Lshortfile)
}
