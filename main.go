// Copyright 2017 Johan Brandhorst. All Rights Reserved.
// See LICENSE for licensing terms.

package main

import (
	"books/server"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	//"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"books/server/proto/library"
)

var logger *logrus.Logger
var host = flag.String("host", "", "host to get LetsEncrypt certificate for")

func init() {
	logger = logrus.StandardLogger()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})
	// Should only be done from init functions
	grpclog.SetLogger(logger)
}

const (
	port = ":9090"
)

func main() {
	lis, err := net.Listen("tcp", port)
	fmt.Println("Listening on", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	library.RegisterBookServiceServer(s, &server.BookService{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
