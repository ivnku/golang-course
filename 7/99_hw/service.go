package main

import (
	"context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strings"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	var acl map[string][]string

	err := json.Unmarshal([]byte(ACLData), &acl)

	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(CheckAccessInterceptor),
	)

	businessService := NewBusinessService(acl)
	adminService := NewAdminService(acl)

	RegisterBizServer(server, businessService)
	RegisterAdminServer(server, adminService)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("Error occurred while starting listening the port", err)
	}

	go func() {
		go server.Serve(listener)
		fmt.Println("Started")
		select {
		case <-ctx.Done():
			fmt.Println("Stopped")
			server.Stop()
		}
	}()

	return nil

}

func checkACL(ctx context.Context, acl map[string][]string) (bool, error) {

	return true, nil
}

func CheckAccessInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Error while getting metadata")
	}

	consumerArr, ok := md["consumer"]

	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Consumer not provided")
	}
	consumer := consumerArr[0]

	var acl map[string][]string
	switch info.Server.(type) {
	case *BusinessService:
		acl = info.Server.(*BusinessService).ACL
	case *AdminService:
		acl = info.Server.(*AdminService).ACL
	}

	canAccess := false
	consumerAllowedMethods := acl[consumer]
	requestedPath := strings.Split(info.FullMethod, "/")
	for _, method := range consumerAllowedMethods {
		aclPath := strings.Split(method, "/")
		if requestedPath[1] == aclPath[1] && (requestedPath[2] == aclPath[2] || aclPath[2] == "*") {
			canAccess = true
			break
		}
	}
	print(requestedPath)

	if !canAccess {
		return nil, status.Errorf(codes.Unauthenticated, "Consumer cannot access method")
	}

	return handler(ctx, req)
}

// Business service
type BusinessService struct {
	ACL map[string][]string
}

func NewBusinessService(acl map[string][]string) *BusinessService {
	return &BusinessService{ACL: acl}
}

func (srv *BusinessService) Check(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (srv *BusinessService) Add(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (srv *BusinessService) Test(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (srv *BusinessService) mustEmbedUnimplementedBizServer() {
	return
}

// Admin service
type AdminService struct {
	ACL map[string][]string
}

func NewAdminService(acl map[string][]string) *AdminService {
	return &AdminService{ACL: acl}
}

func (a AdminService) Logging(nothing *Nothing, server Admin_LoggingServer) error {
	panic("implement me")
}

func (a AdminService) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {
	panic("implement me")
}

func (a AdminService) mustEmbedUnimplementedAdminServer() {
	panic("implement me")
}
