package main

import (
	"context"
	"errors"

	productpb "github.com/amane15/ecom_microservice/proto/product/v1"
)

type ProductGRPCServer struct {
	productpb.UnimplementedProductServiceServer
}

func (s *ProductGRPCServer) CheckVariantExists(ctx context.Context,
	req *productpb.CheckVariantExistsRequest,
) (*productpb.CheckVariantExistsResponse, error) {
	if req.VariantId < 1 {
		return nil, errors.New("invalid variant id")
	}

	exists := req.VariantId == 1

	return &productpb.CheckVariantExistsResponse{
		Exists: exists,
	}, nil
}
