package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/service"
	pb "github.com/zyncc/ecommerce-microservice/services/auth/pkg/types/proto"
	authUtils "github.com/zyncc/ecommerce-microservice/services/auth/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	pb.UnimplementedAuthServiceServer

	log        *zap.Logger
	authSvc    *service.AuthService
	addressSvc *service.AddressService
}

func NewAuthController(log *zap.Logger, authSvc *service.AuthService, addressSvc *service.AddressService) *AuthController {
	return &AuthController{
		log:        log,
		authSvc:    authSvc,
		addressSvc: addressSvc,
	}
}

func (c *AuthController) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	req := dto.SignUpRequest{
		Name:            in.GetName(),
		Email:           in.GetEmail(),
		Password:        in.GetPassword(),
		ConfirmPassword: in.GetConfirmPassword(),
	}

	id, err := c.authSvc.SignUp(ctx, req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.SignUpResponse{
		Id: id.String(),
	}, nil
}

func (c *AuthController) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	req := dto.SignInRequest{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	accessToken, refreshToken, err := c.authSvc.SignIn(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *AuthController) GetSession(ctx context.Context, in *emptypb.Empty) (*pb.GetSessionResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing auth header")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing auth header")
	}

	session, err := authUtils.GetSession(token[0])
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get session")
	}

	return &pb.GetSessionResponse{
		Id:        session.ID.String(),
		Name:      session.Name,
		Email:     session.Email,
		Role:      session.Role,
		CreatedAt: timestamppb.New(session.CreatedAt),
	}, nil
}

func (c *AuthController) RefreshToken(ctx context.Context, in *emptypb.Empty) (*pb.RefreshTokenResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing refresh_token cookie")
	}

	refreshToken := md.Get("refresh_token")
	if len(refreshToken) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing refresh_token cookie")
	}

	accessToken, err := c.authSvc.RefreshToken(ctx, refreshToken[0])
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &pb.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (c *AuthController) CreateAddress(ctx context.Context, in *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	userID, _ := uuid.Parse(in.UserId)
	req := dto.CreateAddressRequest{
		UserID:    userID,
		FirstName: in.GetFirstName(),
		LastName:  in.LastName,
		Email:     in.Email,
		Phone:     in.GetPhone(),
		Address1:  in.GetAddress1(),
		Address2:  in.Address2,
		City:      in.GetCity(),
		State:     in.GetState(),
		Zip:       in.GetZip(),
	}

	addressID, err := c.addressSvc.CreateAddress(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateAddressResponse{
		AddressId: addressID.String(),
	}, nil
}

func (c *AuthController) GetAddressByID(ctx context.Context, in *pb.GetAddressByIDRequest) (*pb.GetAddressByIDResponse, error) {
	addressID, _ := uuid.Parse(in.GetId())

	address, err := c.addressSvc.FindAddressByID(ctx, addressID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetAddressByIDResponse{
		Id:        address.ID.String(),
		UserId:    address.UserID.String(),
		FirstName: address.FirstName,
		LastName:  address.LastName,
		Email:     address.Email,
		Phone:     address.Phone,
		Address1:  address.Address1,
		Address2:  address.Address2,
		City:      address.City,
		State:     address.State,
		Zip:       address.Zip,
		CreatedAt: timestamppb.New(address.CreatedAt),
		UpdatedAt: timestamppb.New(address.UpdatedAt),
	}, nil
}

func (c *AuthController) GetAllAddresses(ctx context.Context, in *pb.GetAllAddressesRequest) (*pb.GetAllAddressesResponse, error) {
	userID, _ := uuid.Parse(in.GetUserId())
	address, err := c.addressSvc.FetchAllAddresses(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var addresses []*pb.GetAddressByIDResponse
	for _, address := range address {
		addressDto := &pb.GetAddressByIDResponse{
			Id:        address.ID.String(),
			UserId:    address.UserID.String(),
			FirstName: address.FirstName,
			LastName:  address.LastName,
			Email:     address.Email,
			Phone:     address.Phone,
			Address1:  address.Address1,
			Address2:  address.Address2,
			City:      address.City,
			State:     address.State,
			Zip:       address.Zip,
			CreatedAt: timestamppb.New(address.CreatedAt),
			UpdatedAt: timestamppb.New(address.UpdatedAt),
		}

		addresses = append(addresses, addressDto)
	}

	return &pb.GetAllAddressesResponse{
		Addresses: addresses,
	}, nil
}
