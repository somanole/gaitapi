// +build !test

// CassandraAccelerationRepo
package accelerationrepo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/constants"
	"github.com/gocql/gocql"
	"code.google.com/p/go-uuid/uuid"
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
	cluster := gocql.NewCluster(constants.CASSANDRA_MASTER)
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraAccelerationRepo) CreateAcceleration(a types.Acceleration) error {
    // insert acceleration
	var err error
	err = nil
	
	sql :=fmt.Sprintf("INSERT INTO accelerations (user_id, timestamp, timestamp_long, x, y, z) VALUES (%v, %v, %v, %v, %v, %v)", a.UserId, a.Timestamp, a.Timestamp, a.X, a.Y, a.Z)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraAccelerationRepo.CreateAcceleration() - Error: %v", err.Error()))
	} 
	
    return err
}

func (repo *CassandraAccelerationRepo) GetAccelerations(userId string) (types.Accelerations, error) {
    // get accelerations
	var accelerations types.Accelerations 
	var err error
	err = nil
	
	var user_id string
	var x, y, z float64 
	var timestamp int64
	
	sql :=fmt.Sprintf("SELECT user_id, x, y, z, timestamp FROM accelerations WHERE user_id = %v", userId)
						
	log.Printf(sql)
	
	iter := session.Query(sql).Iter()
	for iter.Scan(&user_id, &x, &y, &z, &timestamp) {
		accelerations = append(accelerations, types.Acceleration{UserId: uuid.Parse(user_id), X: x, Y: y, Z: z, Timestamp: timestamp})
	}
	
	if err = iter.Close(); err != nil {
		log.Printf(fmt.Sprintf("CassandraAccelerationRepo.GetAccelerations() - Error: %v", err.Error()))
	} 
	
    return accelerations, err
}
