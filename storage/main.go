// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	v1 "backend/gen/board/v1"
	"backend/gen/board/v1/boardv1connect"
	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
)

type StorageServer struct {
	messages []*v1.MessageType
}

func (ps *StorageServer) AddMessage(
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

func (ps *StorageServer) AllMessages(
	ctx context.Context,
	req *connect.Request[v1.AllMessagesRequest],
) (*connect.Response[v1.AllMessagesResponse], error) {
	fmt.Println("Retrieving all the messages")
	res := connect.NewResponse(&v1.AllMessagesResponse{Messages: ps.messages})
	return res, nil
}

func main() {
	mux := http.NewServeMux()
	allMessages := []*v1.MessageType{{
		Id:        1,
		Title:     "First message",
		Message:   "With body",
		Timestamp: time.Now().String(),
	}, {
		Id:        2,
		Title:     "Second message",
		Message:   "Other body",
		Timestamp: time.Now().String(),
	}}
	mux.Handle(boardv1connect.NewMessageServiceHandler(&StorageServer{allMessages}))
	debugReqs := printRequests{h: cors.AllowAll().Handler(mux)}
	addr := "0.0.0.0"
	port := "8081"
	fmt.Printf("Stating server on %s:%s\n", addr, port)
	err := http.ListenAndServe(
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
