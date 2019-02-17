package chefs

import (
	"fmt"
	"net"
	"net/http"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/quocdat32461997/HomeCOOK/api/protos/chefpb"
	"github.com/quocdat32461997/HomeCOOK/internal/cloud"
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

// Converts model.Chef to userpb.Chef
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
		KnownRecepies: chef.KnownRecepies,
	}
}

// Encrypts the chef's password for storage in database
func encrypt(password string) string {
	return password
}

// CreateUser creates new users for the platform
func (s *Server) CreateChef(ctx context.Context, request *chefpb.ChefRequest) (*chefpb.ChefResponse, error) {
	// Extract request data
	data := request.GetChef()
	location := data.GetLocation()

	// Create model and Encrypt password
	user := &models.User{
		FirstName: data.GetFirstName(),
		LastName:  data.GetLastName(),
		Password:  encrypt(data.GetPassword()),
		Location: &models.Location{
			Latitude:  location.GetLatitude(),
			Longitude: location.GetLongitude(),
		},
		Allergens:      data.GetAllergens(),
		FoodPreference: data.GetFoodPreference().String(),
	}

	// Insert into database and add ObjectId to struct
	err := s.Mongo.CreateUser(user)
	fmt.Println(user)

	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		User: convert(user),
	}, nil
}

// GetUser gets a user based on their uuid
func (s *Server) GetUser(ctx context.Context, request *userpb.UserRequest) (*userpb.UserResponse, error) {
	// Extract request data
	uid := request.GetUser().GetId()

	// Pass off uid to database funciotns
	user, err := s.Mongo.GetUser(uid)
	if err != nil {
		return nil, err
	}
	fmt.Print(user)

	return &userpb.UserResponse{
		User: convert(user),
	}, nil
}

// StartUserService starts the user service server
func StartUserService(s *Server) {
	lis, err := net.Listen("tcp", s.Endpoint)
	if err != nil {
		panic(err)
	}

	opts := []grpc.ServerOption{}
	g := grpc.NewServer(opts...)
	userpb.RegisterUserServiceServer(g, s)
	reflection.Register(g)

	if err := g.Serve(lis); err != nil {
		panic(err)
	}
}

// StartUserServiceProxy start's the HTTP to gRPC proxy
func StartUserServiceProxy(s *Server) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := gwruntime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, s.Endpoint, opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
