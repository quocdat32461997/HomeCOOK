package chefs

import (
	"net"
	"net/http"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/quocdat32461997/HomeCOOK/api/protos/chefpb"
	"github.com/quocdat32461997/HomeCOOK/internal/cloud"
	"github.com/quocdat32461997/HomeCOOK/internal/common"
	"github.com/quocdat32461997/HomeCOOK/internal/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the chef service's server
type Server struct {
	Mongo    *cloud.MongoConn
	Endpoint string
	Listener net.Listener
}

// Converts model.Chef to chefpb.Chef
func convert(chef *models.Chef) *chefpb.Chef {

	return &chefpb.Chef{
		Id:        chef.ID.Hex(), // Convert to string for Protocol Buffers model
		FirstName: chef.FirstName,
		LastName:  chef.LastName,
		Password:  chef.Password,
		Location: &chefpb.Location{
			Longitude: chef.Location.Longitude,
			Latitude:  chef.Location.Latitude,
		},
		Rating:        chef.Rating,
		KnownRecipies: chef.KnownRecipies,
	}
}

// CreateChef creates new chefs for the platform
func (s *Server) CreateChef(ctx context.Context, request *chefpb.ChefRequest) (*chefpb.ChefResponse, error) {
	// Extract request data
	data := request.GetChef()
	location := data.GetLocation()
	password, err := common.EncryptPassword(data.GetPassword())
	if err != nil {
		panic(err)
	}

	// Create model and Encrypt password
	chef := &models.Chef{
		FirstName: data.GetFirstName(),
		LastName:  data.GetLastName(),
		Password:  password,
		Location: &models.Location{
			Latitude:  location.GetLatitude(),
			Longitude: location.GetLongitude(),
		},
		Rating:        -1,         // New chefs do not have a raiting yet
		KnownRecipies: []string{}, // New chefs must select what recepies they are familiar with
	}

	// Insert into database and add ObjectId to struct
	err = s.Mongo.CreateChef(chef)

	if err != nil {
		return nil, err
	}

	return &chefpb.ChefResponse{
		Chef: convert(chef),
	}, nil
}

// GetChef gets a chef based on their uuid
func (s *Server) GetChef(ctx context.Context, request *chefpb.ChefRequest) (*chefpb.ChefResponse, error) {
	// Extract request data
	uid := request.GetChef().GetId()

	// Pass off uid to database funciotns
	chef, err := s.Mongo.GetChef(uid)
	if err != nil {
		return nil, err
	}

	return &chefpb.ChefResponse{
		Chef: convert(chef),
	}, nil
}

// StartChefService starts the chef service's server
func StartChefService(s *Server) {
	lis, err := net.Listen("tcp", s.Endpoint)
	if err != nil {
		panic(err)
	}

	opts := []grpc.ServerOption{}
	g := grpc.NewServer(opts...)
	chefpb.RegisterChefServiceServer(g, s)
	reflection.Register(g)

	if err := g.Serve(lis); err != nil {
		panic(err)
	}
}

// StartChefServiceProxy starts the chef service's HTTP to gRPC proxy
func StartChefServiceProxy(s *Server) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := gwruntime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := chefpb.RegisterChefServiceHandlerFromEndpoint(ctx, mux, s.Endpoint, opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}
