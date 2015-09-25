// +build !test

// CassandraRepo
package repo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/acceleration"
	"github.com/gocql/gocql"
)

type (
	CassandraRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraRepo() Repo {
	return &CassandraRepo{}
}

func init() {
	New = NewCassandraRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraRepo) GetAcceleration(userId int64) acceleration.Acceleration {
    //select acceleration by user_id
	a := acceleration.Acceleration{};
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

func (repo *CassandraRepo) GetAllAccelerations() acceleration.Accelerations {
	log.Println("Cassandra - trying to get all accelerations")
	
	// select all accelerations
	var accelerations acceleration.Accelerations
	var user_id, timestamp int64
	var x, y, z float64
	
	iter := session.Query("SELECT * FROM accelerations").Iter()
	for iter.Scan(&user_id, &timestamp, &x, &y, &z) {
		accelerations = append(accelerations, acceleration.Acceleration{UserId: user_id, Timestamp: timestamp, X: x, Y: y, Z: z})
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	
	return accelerations
}

func (repo *CassandraRepo) GetAccelerationsCount() acceleration.AccelerationsCount {
	// select count of all accelerations
	var count int64
	
	if err := session.Query("SELECT COUNT(*) FROM accelerations").Scan(&count); err != nil {
		log.Fatal(err)
	}
	
	response := acceleration.AccelerationsCount{count}
	
	return response
}

func (repo *CassandraRepo) CreateAcceleration(a acceleration.Acceleration) acceleration.Acceleration {
    // insert acceleration
	sql :=fmt.Sprintf("INSERT INTO accelerations (user_id, timestamp, x, y, z) VALUES (%v, %v, %v, %v, %v)", a.UserId, a.Timestamp, a.X, a.Y, a.Z)
	log.Printf(sql)
	if err := session.Query(sql).Exec(); err != nil {
		log.Fatal(err)
	}
	
    return a
}

func (repo *CassandraRepo) DeleteAcceleration(userId int64) error {
    // delete the acceleration from the accelerations table
	var err gocql.Error
	
	sql := fmt.Sprintf("DELETE FROM accelerations WHERE user_id = %v", userId)
	if err := session.Query(sql).Exec(); err != nil {
		log.Fatal(err)
	}
	
	return err
}
