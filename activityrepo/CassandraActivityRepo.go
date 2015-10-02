// CassandraActivityRepo
package activityrepo

import (
	"errors"
	"fmt"
	"log"
	"time"
	"github.com/somanole/gaitapi/types"
	"github.com/gocql/gocql"
	"code.google.com/p/go-uuid/uuid"
)

type (
	CassandraActivityRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraActivityRepo() ActivityRepo {
	return &CassandraActivityRepo{}
}

func init() {
	New = NewCassandraActivityRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraActivityRepo) CreateUserActivity(userId string, ar types.ActivityRequest) error {
    // insert activity in activities
	var a types.Activity
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf(`SELECT email from users_by_id WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		var email string
		if err = session.Query(sql).Scan(&email); err != nil {
				log.Printf(fmt.Sprintf("CreateMessage - Error: %v", err.Error()))
		} else {
			a.DeviceId = ar.DeviceId
			a.DeviceType = ar.DeviceType
			a.IsLoggedIn = ar.IsLoggedIn
			a.UserId = uuid.Parse(userId)
			a.PushToken = ar.PushToken

			if ar.Timestamp != 0 {
				a.Timestamp = int64(time.Unix(ar.Timestamp, 0).UTC().Unix())
			} else {
				a.Timestamp = int64(time.Now().UTC().Unix())
			}
			
			sql := fmt.Sprintf(`INSERT INTO activities (user_id, device_type, 
			device_id, is_logged_in, push_token, timestamp) VALUES (%v, '%v', '%v', %v, '%v', %v)`, 
			a.UserId, a.DeviceType, a.DeviceId, a.IsLoggedIn, a.PushToken, a.Timestamp)
						
			log.Printf(sql)
			if err = session.Query(sql).Exec(); err != nil {
				log.Printf(fmt.Sprintf("CreateActivity - Error: %v", err.Error()))
			} 
		}
	} else {
		log.Printf(fmt.Sprintf("CreateActivity - UserId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
    return err
}

func (repo *CassandraActivityRepo) GetUserActivity(userId string) (types.Activity, error) {
	// get activity for user by id
	log.Printf(fmt.Sprintf("GetUserActivity - Received userId: %v", userId))
	
	var activity types.Activity
	var user_id string
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf(`SELECT user_id, device_type, 
		device_id, is_logged_in, push_token, timestamp 
		FROM activities WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&user_id, &activity.DeviceType, &activity.DeviceId,
		&activity.IsLoggedIn, &activity.PushToken, &activity.Timestamp); err != nil {
				log.Printf(fmt.Sprintf("GetUserActivity - Error: %v", err.Error()))
		} else {
			activity.UserId = uuid.Parse(user_id)
		}
	} else {
		log.Printf(fmt.Sprintf("GetUserActivity - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return activity, err
}
