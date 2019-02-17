package cloud

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/quocdat32461997/HomeCOOK/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	usersCollection = "users"
	chefsCollection = "chefs"
)

// MongoConn provides a connection to the MongoDB Atlas cluster
type MongoConn struct {
	Host       string
	Authorizer string
	Database   string
	Username   string
	Password   string
	Client     *mgo.Database
}

// InitDB initializes the database
func (m *MongoConn) InitDB() (*mgo.Session, error) {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{m.Host},
		Timeout:  10 * time.Second,
		Database: m.Authorizer,
		Username: m.Username,
		Password: m.Password,
	}
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	m.Client = session.DB(m.Database)
	return session, nil
}

// CreateUser creates a user
func (m *MongoConn) CreateUser(user *models.User) error {
	oid := bson.NewObjectId()
	user.ID = oid
	err := m.Client.C(usersCollection).Insert(user)
	if err != nil {
		panic(err)
	}
	return err
}

// GetUser gets a user
func (m *MongoConn) GetUser(id string) (*models.User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Input string is not valid ObjectId hex"),
		)
	}

	oid := bson.ObjectIdHex(id)
	user := &models.User{}
	query := m.Client.C(usersCollection).Find(bson.M{"_id": oid})

	if n, err := query.Count(); n == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Unable to find user with ObjectId %v: %v", id, err),
		)
	}
	query.One(&user)
	return user, nil
}

// CreateChef creates a chef
func (m *MongoConn) CreateChef(chef *models.Chef) error {
	oid := bson.NewObjectId()
	chef.ID = oid
	err := m.Client.C(chefsCollection).Insert(chef)
	if err != nil {
		panic(err)
	}
	return err
}

// GetChef gets a chef
func (m *MongoConn) GetChef(id string) (*models.Chef, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Input string is not valid ObjectId hex"),
		)
	}

	oid := bson.ObjectIdHex(id)
	chef := &models.Chef{}
	query := m.Client.C(chefsCollection).Find(bson.M{"_id": oid})

	if n, err := query.Count(); n == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Unable to find chef with ObjectId %v: %v", id, err),
		)
	}
	query.One(&chef)
	return chef, nil
}
