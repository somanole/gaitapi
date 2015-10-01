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

func (repo *CassandraAccelerationRepo) GetAcceleration(userId int64) types.Acceleration {
    //select acceleration by user_id
	a := types.Acceleration{};
	var user_id, timestamp int64
	var x, y, z float64
	
	sql := fmt.Sprintf("SELECT * FROM accelerations WHERE user_id=%v", userId)
	iter := session.Query(sql).Iter()
	for iter.Scan(&user_id, &timestamp, &x, &y, &z) {
		a.UserId = user_id
		a.Timestamp = timestamp
		a.X = x
		a.Y = y
		a.Z = z
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	
	return a
}

func (repo *CassandraAccelerationRepo) GetAllAccelerations() types.Accelerations {
	log.Println("Cassandra - trying to get all accelerations")
	
	// select all accelerations
	var accelerations types.Accelerations
	var user_id, timestamp int64
	var x, y, z float64
	
	iter := session.Query("SELECT * FROM accelerations").Iter()
	for iter.Scan(&user_id, &timestamp, &x, &y, &z) {
		accelerations = append(accelerations, types.Acceleration{UserId: user_id, Timestamp: timestamp, X: x, Y: y, Z: z})
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	
	return accelerations
}

func (repo *CassandraAccelerationRepo) GetAccelerationsCount() types.AccelerationsCount {
	// select count of all accelerations
	var count int64
	
	if err := session.Query("SELECT COUNT(*) FROM accelerations").Scan(&count); err != nil {
		log.Fatal(err)
	}
	
	response := types.AccelerationsCount{count}
	
	return response
}

func (repo *CassandraAccelerationRepo) CreateAcceleration(a types.Acceleration) types.Acceleration {
    // insert acceleration
	sql :=fmt.Sprintf("INSERT INTO accelerations (user_id, timestamp, x, y, z) VALUES (%v, %v, %v, %v, %v)", a.UserId, a.Timestamp, a.X, a.Y, a.Z)
	log.Printf(sql)
	if err := session.Query(sql).Exec(); err != nil {
		log.Fatal(err)
	}
	
    return a
}

func (repo *CassandraAccelerationRepo) DeleteAcceleration(userId int64) error {
    // delete the acceleration from the accelerations table
	var err gocql.Error
	
	sql := fmt.Sprintf("DELETE FROM accelerations WHERE user_id = %v", userId)
	if err := session.Query(sql).Exec(); err != nil {
		log.Fatal(err)
	}
	
	return err
}
