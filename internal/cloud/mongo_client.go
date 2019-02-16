package cloud

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/globalsign/mgo"
	"github.com/quocdat32461997/HomeCOOK/internal/models"
	// "github.com/globalsign/mgo/bson"
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
func (m *MongoConn) InitDB() *mgo.Session {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{m.Host},
		Timeout:  10 * time.Second,
		Database: m.Database,
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
		panic(nil)
	}
	session.SetMode(mgo.Monotonic, true)
	m.Client = session.DB(m.Database)
	return session
}

// CreateUser creates a user
func (m *MongoConn) CreateUser(user *models.User) {
	fmt.Println("Testing")
	m.Client.C("user").Insert(testUser)

}
