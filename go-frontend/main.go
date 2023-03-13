// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	v1 "backend/gen/board/v1"
	"backend/gen/board/v1/boardv1connect"
	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
)

func main() {
	storageAddress := os.Getenv("BACKEND_ADDRESS")
	if storageAddress == "" {
		storageAddress = "localhost:8080"
	}
	client := boardv1connect.NewMessageServiceClient(http.DefaultClient,
		fmt.Sprintf("http://%s", storageAddress))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		resp, err := client.AllMessages(context.Background(), connect.NewRequest(&v1.AllMessagesRequest{}))
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(fmt.Sprintf("error from the backend: %s", err.Error())))
		}
		if resp == nil {
			writer.WriteHeader(500)
			writer.Write([]byte(fmt.Sprintf("did not get any response")))
		}
		for _, message := range resp.Msg.Messages {
			fmt.Println("Got message with title: %s\n", message.Title)
		}
		reply, err := json.Marshal(resp.Msg.Messages)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(fmt.Sprintf("failed to marshal: %s", err.Error())))
		}
		writer.Write(reply)
		writer.WriteHeader(200)
	})
	debugReqs := printRequests{h: cors.AllowAll().Handler(mux)}
	addr := "0.0.0.0"
	port := "80"

	_, err := client.AllMessages(context.Background(), connect.NewRequest(&v1.AllMessagesRequest{}))
	if err == nil {
		fmt.Println("Was able to reach the backend")
	} else {
		fmt.Printf("Unable to reach backend with: %s\n", err.Error())
	}

	fmt.Printf("Stating server on %s:%s\n", addr, port)
	err = http.ListenAndServe(
		fmt.Sprintf("%s:%s", addr, port),
		debugReqs,
	)
	log.Fatalf("listen failed: %v", err)
}

type printRequests struct{ h http.Handler }

func (p printRequests) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("%s %s\n", request.Method, request.RequestURI)

	p.h.ServeHTTP(writer, request)
}
