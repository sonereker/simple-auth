package client

import (
	"context"
	"fmt"
	"github.com/sonereker/simple-auth/pb/v1"
	"google.golang.org/grpc"
)

//NewUserClientWithConnection dials a connection to gRPC server and initializes a new UserClient using the connection
func NewUserClientWithConnection(grpcServerAddr string, attachToken bool, loginRequest *pb.LoginRequest) (*grpc.ClientConn, pb.UserClient) {
	var uc pb.UserClient
	conn, _ := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	uc = pb.NewUserClient(conn)

	if attachToken {
		login, err := uc.Login(context.Background(), loginRequest)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		interceptor := AuthInterceptor{
			Token: login.Token,
		}
		conn, _ = grpc.Dial(
			grpcServerAddr,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(interceptor.GetUnaryInterceptor()),
		)
		uc = pb.NewUserClient(conn)
	}
	return conn, uc
}
