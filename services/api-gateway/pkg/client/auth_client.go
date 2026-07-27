package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/zyncc/ecommerce-microservice/services/auth/pkg/types/proto"

	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"
	"go.uber.org/zap"
)

type AuthClient interface {
	SignUp(ctx context.Context, req *dto.SignUpRequest) (string, error)
	SignIn(ctx context.Context, req *dto.SignInRequest) (dto.SignInResponse, error)
	GetSession(ctx context.Context, r *http.Request) (types.Session, error)
	RefreshToken(ctx context.Context, r *http.Request) (string, error)
	CreateAddress(ctx context.Context, req dto.CreateAddressRequest) (uuid.UUID, error)
	GetAddressByID(ctx context.Context, id uuid.UUID) (dto.AddressResponse, error)
	GetAllAddresses(ctx context.Context, userID uuid.UUID) ([]dto.AddressResponse, error)
}

type GRPCAuthClient struct {
	client pb.AuthServiceClient
	log    *zap.Logger
	conn   *grpc.ClientConn
}

func NewGRPCAuthClient(log *zap.Logger, addr string) (*GRPCAuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCAuthClient{
		log:    log,
		conn:   conn,
		client: pb.NewAuthServiceClient(conn),
	}, nil
}

func (c *GRPCAuthClient) SignUp(ctx context.Context, req *dto.SignUpRequest) (string, error) {
	in := &pb.SignUpRequest{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	resp, err := c.client.SignUp(ctx, in)
	if err != nil {
		c.log.Error("auth service returned error", zap.Error(err))
		return "", err
	}

	return resp.GetId(), nil
}

func (c *GRPCAuthClient) SignIn(ctx context.Context, req *dto.SignInRequest) (dto.SignInResponse, error) {
	in := &pb.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := c.client.SignIn(ctx, in)
	if err != nil {
		c.log.Error("auth service returned error", zap.Error(err))
		return dto.SignInResponse{}, err
	}

	dto := dto.SignInResponse{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
	}

	return dto, nil
}

func (c *GRPCAuthClient) GetSession(ctx context.Context, r *http.Request) (types.Session, error) {
	token, err := utils.ExtractAuthHeader(r)
	if err != nil {
		return types.Session{}, err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	resp, err := c.client.GetSession(ctx, &emptypb.Empty{})
	if err != nil {
		return types.Session{}, err
	}

	id, _ := uuid.Parse(resp.GetId())
	dto := types.Session{
		ID:        id,
		Name:      resp.GetName(),
		Email:     resp.GetEmail(),
		Role:      resp.GetRole(),
		CreatedAt: resp.GetCreatedAt().AsTime(),
	}

	return dto, nil
}

func (c *GRPCAuthClient) RefreshToken(ctx context.Context, r *http.Request) (string, error) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", errors.New("refresh_token cookie is missing")
	}
	refreshToken := refreshTokenCookie.Value

	ctx = metadata.AppendToOutgoingContext(ctx, "refresh_token", refreshToken)
	resp, err := c.client.RefreshToken(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return resp.GetAccessToken(), nil
}

func (c *GRPCAuthClient) CreateAddress(ctx context.Context, req dto.CreateAddressRequest) (uuid.UUID, error) {
	in := &pb.CreateAddressRequest{
		UserId:    req.UserID.String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address1:  req.Address1,
		Address2:  req.Address2,
		City:      req.City,
		State:     req.State,
		Zip:       req.Zip,
	}

	resp, err := c.client.CreateAddress(ctx, in)
	if err != nil {
		return uuid.Nil, err
	}

	addressID, err := uuid.Parse(resp.GetAddressId())
	if err != nil {
		return uuid.Nil, err
	}

	return addressID, nil
}

func (c *GRPCAuthClient) GetAddressByID(ctx context.Context, id uuid.UUID) (dto.AddressResponse, error) {
	in := &pb.GetAddressByIDRequest{
		Id: id.String(),
	}

	resp, err := c.client.GetAddressByID(ctx, in)
	if err != nil {
		return dto.AddressResponse{}, err
	}

	addressID, _ := uuid.Parse(resp.GetId())
	userID, _ := uuid.Parse(resp.GetUserId())

	dto := dto.AddressResponse{
		ID:        addressID,
		UserID:    userID,
		FirstName: resp.GetFirstName(),
		LastName:  resp.LastName,
		Email:     resp.GetEmail(),
		Phone:     resp.GetPhone(),
		Address1:  resp.GetAddress1(),
		Address2:  resp.Address2,
		City:      resp.GetCity(),
		State:     resp.GetState(),
		Zip:       resp.GetZip(),
		CreatedAt: resp.GetCreatedAt().AsTime(),
		UpdatedAt: resp.GetUpdatedAt().AsTime(),
	}

	return dto, nil
}

func (c *GRPCAuthClient) GetAllAddresses(ctx context.Context, userID uuid.UUID) ([]dto.AddressResponse, error) {
	in := &pb.GetAllAddressesRequest{
		UserId: userID.String(),
	}

	resp, err := c.client.GetAllAddresses(ctx, in)
	if err != nil {
		return nil, err
	}

	var addresses []dto.AddressResponse

	for _, address := range resp.GetAddresses() {
		id, _ := uuid.Parse(address.GetId())
		userID, _ := uuid.Parse(address.GetUserId())

		dto := dto.AddressResponse{
			ID:        id,
			UserID:    userID,
			FirstName: address.GetFirstName(),
			LastName:  address.LastName,
			Email:     address.GetEmail(),
			Phone:     address.GetPhone(),
			Address1:  address.GetAddress1(),
			Address2:  address.Address2,
			City:      address.GetCity(),
			State:     address.GetState(),
			Zip:       address.GetZip(),
			CreatedAt: address.GetCreatedAt().AsTime(),
			UpdatedAt: address.GetUpdatedAt().AsTime(),
		}

		addresses = append(addresses, dto)
	}

	return addresses, nil
}
