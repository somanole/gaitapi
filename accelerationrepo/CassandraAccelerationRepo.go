// +build !test

// CassandraAccelerationRepo
package accelerationrepo

import (
	"fmt"
	"log"
	"time"
	"errors"
	"github.com/somanole/gaitapi/types"
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
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraAccelerationRepo) CreateAcceleration(userId string, ar types.AccelerationRequest) error {
    // insert acceleration
	var a types.Acceleration
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf(`SELECT email from users_by_id WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		var email string
		if err = session.Query(sql).Scan(&email); err != nil {
				log.Printf(fmt.Sprintf("CreateMessage - Error: %v", err.Error()))
		} else {
			a.X = ar.X
			a.Y = ar.Y
			a.Z = ar.Z
			a.UserId = uuid.Parse(userId)

			if ar.Timestamp != 0 {
				a.Timestamp = int64(time.Unix(ar.Timestamp, 0).UTC().Unix())
			} else {
				a.Timestamp = int64(time.Now().UTC().Unix())
			}
			
			sql :=fmt.Sprintf("INSERT INTO accelerations (user_id, timestamp, x, y, z) VALUES (%v, %v, %v, %v, %v)", a.UserId, a.Timestamp, a.X, a.Y, a.Z)
						
			log.Printf(sql)
			if err = session.Query(sql).Exec(); err != nil {
				log.Printf(fmt.Sprintf("CreateAcceleration - Error: %v", err.Error()))
			} 
		}
	} else {
		log.Printf(fmt.Sprintf("CreateAcceleration - UserId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
    return err
}
