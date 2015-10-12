// +build !test

// CassandraAccelerationRepo
package accelerationrepo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
	"github.com/gocql/gocql"
)

type (
	CassandraAccelerationRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraAccelerationRepo() AccelerationRepo {
	return &CassandraAccelerationRepo{}
}

func init() {
	New = NewCassandraAccelerationRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraAccelerationRepo) CreateAcceleration(a types.Acceleration) error {
    // insert acceleration
	var err error
	err = nil
	
	sql :=fmt.Sprintf("INSERT INTO accelerations (user_id, timestamp, x, y, z) VALUES (%v, %v, %v, %v, %v)", a.UserId, a.Timestamp, a.X, a.Y, a.Z)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraAccelerationRepo.CreateAcceleration() - Error: %v", err.Error()))
	} 
	
    return err
}
