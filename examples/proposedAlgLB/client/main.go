/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	ka "google.golang.org/grpc/keepalive"
	"log"
	"time"
)

const (
	defaultName = "proposedLB"
)

var (
	addr   = flag.String("addr", "localhost:50051", "the address to connect to")
	name   = flag.String("name", defaultName, "Name to greet")
	params = ka.ClientParameters{
		Time:                3 * time.Millisecond,
		Timeout:             6 * time.Millisecond,
		PermitWithoutStream: true,
	}
)

func main() {
	fmt.Print(time.Now().String())
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
   "loadBalancingConfig":[
      {
         "round_robin":{
           "childPolicy":[
               {
                  "proposedLB":{}
               }
            ]
         }
      }
   ]
}`), grpc.WithKeepaliveParams(params),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Pass the struct to SetKeepAliveParams()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	fmt.Print(time.Now().String())
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
