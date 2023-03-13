// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	v1 "backend/gen/board/v1"
	"backend/gen/board/v1/boardv1connect"
	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
)

type MessageServer struct {
	client boardv1connect.MessageServiceClient
}

func (ps *MessageServer) AddMessage(
	ctx context.Context,
	req *connect.Request[v1.AddMessageRequest],
) (*connect.Response[v1.AddMessageResponse], error) {
	fmt.Println("Adding a new message")
	// Talk to the other service
	// Reply
	res := connect.NewResponse(&v1.AddMessageResponse{
		ResultCode: 1,
		Result: &v1.AddMessageResponse_Message{
			Message: &v1.MessageType{
				Id:        0,
				Title:     req.Msg.Message.Title,
				Message:   req.Msg.Message.Message,
				Timestamp: req.Msg.Message.Timestamp,
			},
		},
	})
	return res, nil
}

func (ps *MessageServer) AllMessages(
	ctx context.Context,
	req *connect.Request[v1.AllMessagesRequest],
) (*connect.Response[v1.AllMessagesResponse], error) {
	fmt.Println("Retrieving all the messages")
	messages, err := ps.client.AllMessages(ctx, connect.NewRequest(&v1.AllMessagesRequest{}))
	if err != nil {
		return nil, err
	}
	fmt.Printf("result: \n")
	for _, message := range messages.Msg.Messages {
		fmt.Printf("%d: %s '%s'\n", message.Id, message.Title, message.Message)
	}

	res := connect.NewResponse(&v1.AllMessagesResponse{Messages: messages.Msg.Messages})
	return res, nil
}

func main() {
	storageAddress := os.Getenv("STORAGE_ADDRESS")
	if storageAddress == "" {
		storageAddress = "localhost:8081"
	}
	client := boardv1connect.NewMessageServiceClient(http.DefaultClient,
		fmt.Sprintf("http://%s", storageAddress))
	mux := http.NewServeMux()
	// The generated constructors return a path and a plain net/http
	// handler.
	svc := &MessageServer{client}
	mux.Handle(boardv1connect.NewMessageServiceHandler(svc))
	debugReqs := printRequests{h: cors.AllowAll().Handler(mux)}
	addr := "0.0.0.0"

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r, err := svc.AllMessages(context.Background(), nil)
	if err != nil {
		fmt.Printf("error talking with storage: %s", err)
	} else if r == nil {
		fmt.Println("did not get any reply")
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
