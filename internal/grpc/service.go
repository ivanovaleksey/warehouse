package grpc

import (
	"context"

	"warehouse/api/warehousepb"
	"warehouse/internal/services/products"
)

type Service struct {
	warehousepb.UnimplementedWarehouseServiceServer
	productsSrv products.Service
}

func NewService(productsSrv products.Service) *Service {
	return &Service{
		productsSrv: productsSrv,
	}
}

func (srv *Service) GetProducts(ctx context.Context, _ *warehousepb.GetProductsRequest) (*warehousepb.GetProductsResponse, error) {
	prodsWithStock, err := srv.productsSrv.GetProductsWithStock(ctx)
	if err != nil {
		return nil, err
	}

	resp := &warehousepb.GetProductsResponse{
		Items: make([]*warehousepb.Product, len(prodsWithStock)),
	}
	for _, prod := range prodsWithStock {
		item := &warehousepb.Product{
			Id:    prod.ID,
			Name:  prod.Name,
			Price: prod.Price,
			Stock: prod.Stock,
		}
		item.Articles = make([]*warehousepb.Product_Article, 0, len(prod.Articles))
		for _, art := range prod.Articles {
			item.Articles = append(item.Articles, &warehousepb.Product_Article{
				Id:       art.ID,
				Quantity: art.Quantity,
			})
		}
		resp.Items = append(resp.Items, item)
	}

	return resp, nil
}

func (srv *Service) RemoveProduct(ctx context.Context, req *warehousepb.RemoveProductRequest) (*warehousepb.RemoveProductResponse, error) {
	err := srv.productsSrv.RemoveProduct(ctx, req.Id, 1)
	if err != nil {
		return nil, err
	}
	return &warehousepb.RemoveProductResponse{}, nil
}
