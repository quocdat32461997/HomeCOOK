package users

import (
	"fmt"
	"net"
	"net/http"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/quocdat32461997/HomeCOOK/api/protos/userpb"
	"github.com/quocdat32461997/HomeCOOK/internal/cloud"
	"github.com/quocdat32461997/HomeCOOK/internal/models"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/nickelapp/nickel-api/api/protos/contentpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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
		Id:        user.ID,
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

// Encrypts the users password for storage in database
func encrypt(password string) string {
	return password
}

// CreateUser creates new users for the platform
func (s *Server) CreateUser(ctx context.Context, request *userpb.UserRequest) (*userpb.UserResponse, error) {
	// Extract request data
	data := request.GetUser()
	location := data.GetLocation()

	// Generate UUID
	id := uuid.Must(uuid.NewV4())

	// Create model and Encrypt password to model
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

	// Insert into database
	err := s.Mongo.CreateUser(user)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unable to add user to database"),
		)
	}

	return &userpb.UserResponse{
		User: convert(user),
	}, nil
}

/*
func (s *Server) CreateUser(ctx context.Context, request *userpb.UserRequest) (*userpb.UserResponse, error) {
	// Extract request data
	uid := request.GetUser().GetUserId()

	// Insert into database
	err := s.Mongo.Client.FindUser(uid)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Unable to find user with id: %v", err),
		)
	}

	return &userpb.UserResponse{
		User: convert(user),
	}, nil
}
*/

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

	if err := g.Serve(s.Listener); err != nil {
		panic(err)
	}
}

// StartUserServiceProxy start's the HTTP to gRPC proxy
func StartUserServiceProxy(s *Server) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := gwruntime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := contentpb.RegisterContentServiceHandlerFromEndpoint(ctx, mux, s.Endpoint, opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
