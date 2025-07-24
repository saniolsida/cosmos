package keeper

import (
	"context"

	mymodulev1beta1 "github.com/cosmos/cosmos-sdk/x/mymodule/types/cosmos/mymodule/v1beta1"
)

type queryServer struct {
	Keeper
	mymodulev1beta1.UnimplementedQueryServer
}

func NewQueryServer(k Keeper) mymodulev1beta1.QueryServer {
	return &queryServer{Keeper: k}
}

func (q queryServer) Hello(ctx context.Context, req *mymodulev1beta1.QueryHelloRequest) (*mymodulev1beta1.QueryHelloResponse, error) {
	return &mymodulev1beta1.QueryHelloResponse{
		Message: "Hello, " + req.Name + "!",
	}, nil
}
