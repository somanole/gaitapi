// CassandraActivityRepo
package activityrepo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/constants"
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
	cluster := gocql.NewCluster(constants.CASSANDRA_MASTER)
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraActivityRepo) CreateUserActivity(a types.Activity) error {
    // insert activity in activities
	var err error
	err = nil
	
	sql := fmt.Sprintf(`INSERT INTO activities (user_id, device_type, 
	device_id, is_logged_in, push_token, timestamp) VALUES (%v, '%v', '%v', %v, '%v', %v)`, 
	a.UserId, a.DeviceType, a.DeviceId, a.IsLoggedIn, a.PushToken, a.Timestamp)
						
	log.Printf(sql)
	
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraActivityRepo.CreateActivity() - Error: %v", err.Error()))
	} 

    return err
}

func (repo *CassandraActivityRepo) GetUserActivity(userId string) (types.Activity, error) {
	// get activity for user by id
	log.Printf(fmt.Sprintf("CassandraActivityRepo.GetUserActivity() - Received userId: %v", userId))
	
	var activity types.Activity
	var user_id string
	var err error
	err = nil
	
	sql := fmt.Sprintf(`SELECT user_id, device_type, 
	device_id, is_logged_in, push_token, timestamp 
	FROM activities WHERE user_id = %v LIMIT 1`, userId)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&user_id, &activity.DeviceType, &activity.DeviceId,
	&activity.IsLoggedIn, &activity.PushToken, &activity.Timestamp); err != nil {
		log.Printf(fmt.Sprintf("CassandraActivityRepo.GetUserActivity() - Error: %v", err.Error()))
	} else {
		activity.UserId = uuid.Parse(user_id)
	}
	
	return activity, err
}
