package users

import (
	"net"
	"net/http"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/quocdat32461997/HomeCOOK/api/protos/userpb"
	"github.com/quocdat32461997/HomeCOOK/internal/cloud"
	"github.com/quocdat32461997/HomeCOOK/internal/common"
	"github.com/quocdat32461997/HomeCOOK/internal/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the user service's server
type Server struct {
	Mongo    *cloud.MongoConn
	Endpoint string
	Listener net.Listener
}

// Converts model.User to userpb.User
func convert(user *models.User) *userpb.User {
	var foodPreference userpb.User_FoodPreference
	switch user.FoodPreference {
	case "NONE":
		foodPreference = userpb.User_NONE
	case "VEGAN":
		foodPreference = userpb.User_VEGAN
	case "VEGETARIAN":
		foodPreference = userpb.User_VEGETARIAN
	case "PESCETARIAN":
		foodPreference = userpb.User_PESCETARIAN
	}

	return &userpb.User{
		Id:        user.ID.Hex(), // Convert to string for Protocol Buffers model
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Location: &userpb.Location{
			Longitude: user.Location.Longitude,
			Latitude:  user.Location.Latitude,
		},
		Allergens:      user.Allergens,
		FoodPreference: foodPreference,
	}
}

// CreateUser creates new users for the platform
func (s *Server) CreateUser(ctx context.Context, request *userpb.UserRequest) (*userpb.UserResponse, error) {
	// Extract request data
	data := request.GetUser()
	location := data.GetLocation()
	password, err := common.EncryptPassword(data.GetPassword())
	if err != nil {
		panic(err)
	}

	// Create model and Encrypt password
	user := &models.User{
		FirstName: data.GetFirstName(),
		LastName:  data.GetLastName(),
		Password:  password,
		Location: &models.Location{
			Latitude:  location.GetLatitude(),
			Longitude: location.GetLongitude(),
		},
		Allergens:      data.GetAllergens(),
		FoodPreference: data.GetFoodPreference().String(),
	}

	// Insert into database and add ObjectId to struct
	err = s.Mongo.CreateUser(user)

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
