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
	"sync"
	"time"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	var acl map[string][]string

	err := json.Unmarshal([]byte(ACLData), &acl)

	if err != nil {
		return err
	}

	businessService := NewBusinessService(acl)
	adminService := NewAdminService(acl)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(CheckAccessInterceptor, LogInterceptor(adminService)),
		grpc.ChainStreamInterceptor(CheckAccessStreamInterceptor, LogStreamInterceptor(adminService)),
	)

	RegisterBizServer(server, businessService)
	RegisterAdminServer(server, adminService)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("Error occurred while starting listening the port", err)
	}

	go adminService.watchForLogs(ctx)
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

/**
 * @Description: Get consumer name from context
 * @param ctx
 * @return string
 * @return error
 */
func getConsumer(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "", status.Errorf(codes.DataLoss, "Error while getting metadata")
	}

	consumerArr, ok := md["consumer"]

	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "Consumer not provided")
	}

	return consumerArr[0], nil
}

/**
 * @Description: Check if the user can access specific method
 * @param ctx
 * @param acl
 * @param fullMethod
 * @return error
 */
func CheckACL(ctx context.Context, acl map[string][]string, fullMethod string) error {
	consumer, err := getConsumer(ctx)

	if err != nil {
		return err
	}

	canAccess := false
	consumerAllowedMethods := acl[consumer]
	requestedPath := strings.Split(fullMethod, "/")
	for _, method := range consumerAllowedMethods {
		aclPath := strings.Split(method, "/")
		if requestedPath[1] == aclPath[1] && (requestedPath[2] == aclPath[2] || aclPath[2] == "*") {
			canAccess = true
			break
		}
	}
	print(requestedPath)

	if !canAccess {
		return status.Errorf(codes.Unauthenticated, "Consumer cannot access method")
	}

	return nil
}

/**
 * @Description: Check access before executing the functions
 * @return func
 */
func CheckAccessInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	var acl map[string][]string
	switch info.Server.(type) {
	case *BusinessService:
		acl = info.Server.(*BusinessService).ACL
	case *AdminService:
		acl = info.Server.(*AdminService).ACL
	}

	err := CheckACL(ctx, acl, info.FullMethod)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Consumer cannot access method")
	}

	return handler(ctx, req)
}

/**
 * @Description: Check access before executing the functions (Stream)
 * @return func
 */
func CheckAccessStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	var acl map[string][]string
	switch srv.(type) {
	case *BusinessService:
		acl = srv.(*BusinessService).ACL
	case *AdminService:
		acl = srv.(*AdminService).ACL
	}

	err := CheckACL(ss.Context(), acl, info.FullMethod)

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Consumer cannot access method")
	}

	return handler(srv, ss)
}

/**
 * @Description: Perform logging before accessing the method
 * @param adm
 * @return grpc.UnaryServerInterceptor
 */
func LogInterceptor(adm *AdminService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := adm.Log(ctx, info.FullMethod)

		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

/**
 * @Description: Perform logging before accessing the method (Stream)
 * @param adm
 * @return grpc.UnaryServerInterceptor
 */
func LogStreamInterceptor(adm *AdminService) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler, ) error {
		err := adm.Log(ss.Context(), info.FullMethod)

		if err != nil {
			return err
		}

		return handler(srv, ss)
	}
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
	ACL           map[string][]string
	Logs          chan *Event
	LogSubs       []chan *Event
	LogSubsMutex  sync.RWMutex
	StatSubs      []chan *Event
	StatSubsMutex sync.RWMutex
}

func NewAdminService(acl map[string][]string) *AdminService {
	return &AdminService{
		ACL:      acl,
		Logs:     make(chan *Event),
		LogSubs:  make([]chan *Event, 0),
		StatSubs: make([]chan *Event, 0),
	}
}

func (a *AdminService) Logging(nothing *Nothing, server Admin_LoggingServer) error {
	ch := make(chan *Event)

	a.LogSubsMutex.Lock()
	a.LogSubs = append(a.LogSubs, ch)
	a.LogSubsMutex.Unlock()

	//go func() {
	for {
		select {
		case logEvent := <-ch:
			err := server.Send(logEvent)
			if err != nil {
				a.LogSubsMutex.Lock()
				a.LogSubs = removeChannel(a.LogSubs, ch)
				a.LogSubsMutex.Unlock()
				//break
				return err
			}
		case <-server.Context().Done():
			a.LogSubsMutex.Lock()
			a.LogSubs = removeChannel(a.LogSubs, ch)
			a.LogSubsMutex.Unlock()
			//break
			return server.Context().Err()
		}
	}
	//}()

	//return nil
}

/**
 * @Description: Remove the specific channel from array
 * @param slice
 * @param element
 * @return []chan
 */
func removeChannel(slice []chan *Event, element chan *Event) []chan *Event {
	for index, ch := range slice {
		if ch == element {
			return append(slice[:index], slice[index+1:]...)
		}
	}
	return slice
}

/**
 * @Description: Log the Event to "Logs" channel from which
 * it will be distributed to other channels
 * @receiver a
 * @param ctx
 * @param fullMethod
 * @return error
 */
func (a *AdminService) Log(ctx context.Context, fullMethod string) error {
	consumer, err := getConsumer(ctx)

	if err != nil {
		return err
	}

	a.Logs <- &Event{
		Timestamp: 0,
		Consumer:  consumer,
		Method:    fullMethod,
		Host:      "127.0.0.1:",
	}

	return nil
}

/**
 * @Description: Watch for every write to Logs channel and push the event
 * to log and stat subscribers
 * @receiver a
 * @param ctx
 */
func (a *AdminService) watchForLogs(ctx context.Context) {
	for {
		select {
		case logEvent := <-a.Logs:
			for _, ls := range a.LogSubs {
				ls <- logEvent
			}
			for _, ss := range a.StatSubs {
				ss <- logEvent
			}
		case <-ctx.Done():
			break
		}
	}
}

/**
 * @Description: Collect statistic and send to client with specific time interval
 * @receiver a
 * @param interval
 * @param server
 * @return error
 */
func (a *AdminService) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {
	ticker := time.NewTicker(time.Duration(interval.GetIntervalSeconds()) * time.Second)

	ch := make(chan *Event)

	a.StatSubsMutex.Lock()
	a.StatSubs = append(a.StatSubs, ch)
	a.StatSubsMutex.Unlock()

	statData := createEmptyStatData()

	for {
		select {
		case stat := <-ch:
			statData.ByConsumer[stat.Consumer]++
			statData.ByMethod[stat.Method]++
		case <-ticker.C:
			err := server.Send(statData)
			if err != nil {
				a.StatSubsMutex.Lock()
				a.StatSubs = removeChannel(a.StatSubs, ch)
				a.StatSubsMutex.Unlock()

				ticker.Stop()
				return err
			}
			statData = createEmptyStatData()
		case <-server.Context().Done():
			a.StatSubsMutex.Lock()
			a.StatSubs = removeChannel(a.StatSubs, ch)
			a.StatSubsMutex.Unlock()

			ticker.Stop()
			return server.Context().Err()
		}
	}
}

/**
 * @Description: Create empty Stat data
 * @return *Stat
 */
func createEmptyStatData() *Stat {
	return &Stat{ByMethod: make(map[string]uint64), ByConsumer: make(map[string]uint64), Timestamp: time.Now().Unix()}
}

func (a *AdminService) mustEmbedUnimplementedAdminServer() {
	panic("implement me")
}
