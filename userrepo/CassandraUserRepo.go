// +build !test

// CassandraUserRepo
package userrepo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/constants"
	"github.com/gocql/gocql"
	"code.google.com/p/go-uuid/uuid"
)

type (
	CassandraUserRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraUserRepo() UserRepo {
	return &CassandraUserRepo{}
}

func init() {
	New = NewCassandraUserRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster(constants.CASSANDRA_MASTER)
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraUserRepo) CreateUser(u types.User) (types.CreateUserResponse, error) {
    // insert user in users
	var response types.CreateUserResponse
	var err error
	var sql string
	err = nil
	
	sql = fmt.Sprintf(`INSERT INTO users (user_id, username, 
	facebook_access_token, twitter_access_token, google_access_token, 
	push_token, device_type, email, password, is_test,
	is_anonymous, gender_preference, timestamp) 
	VALUES (%v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, %v,
	'%v', %v)`, 
	u.UserId, u.Username, u.FacebookAccessToken, u.TwitterAccessToken, 
	u.GoogleAccessToken, u.PushToken, u.DeviceType, 
	u.Email, u.Password, u.IsTest, u.IsAnonymous, u.GenderPreference, 
	u.Timestamp)
				
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraUserRepo.CreateUser() - Error: %v", err.Error()))
	} else {
		response.UserId = u.UserId
		response.Username = u.Username
		response.Timestamp = u.Timestamp
			
		if (!u.IsAnonymous)	{
			sql = fmt.Sprintf("INSERT INTO users_by_email (email, user_id) VALUES ('%v', %v)", u.Email, u.UserId)
				
			log.Printf(sql)
			if err = session.Query(sql).Exec(); err != nil {
				log.Printf(fmt.Sprintf("CreateUser - Error: %v", err.Error()))
			} 
		}	
				
		sql = fmt.Sprintf("INSERT INTO users_extra_info (user_id, walking_progress, timestamp) VALUES (%v, %v, %v)", u.UserId, 0, u.Timestamp)
					
		log.Printf(sql)
		if err = session.Query(sql).Exec(); err != nil {
			log.Printf(fmt.Sprintf("CassandraUserRepo.CreateUser() - Error: %v", err.Error()))
		} 
	}
	
    return response, err
}

func (repo *CassandraUserRepo) GetUser(userId string) (types.User, error) {
	// get user by id
	var user types.User
	var user_id string
	var err error
	err = nil
	
	sql := fmt.Sprintf(`SELECT user_id, username, facebook_access_token, 
	twitter_access_token, google_access_token, push_token, 
	device_type, email, password, is_test, is_anonymous, gender_preference, 
	timestamp FROM users WHERE user_id = %v LIMIT 1`, userId)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&user_id, 
		&user.Username, &user.FacebookAccessToken, &user.TwitterAccessToken, &user.GoogleAccessToken,
		&user.PushToken, &user.DeviceType, &user.Email,
		&user.Password, &user.IsTest, &user.IsAnonymous, &user.GenderPreference, &user.Timestamp); err != nil {
			log.Printf(fmt.Sprintf("CassandraUserRepo.GetUser() - Error: %v", err.Error()))
	} else {
		user.UserId = uuid.Parse(user_id)
	}
	
	return user, err
}

func (repo *CassandraUserRepo) UpdateUser(user types.User) (types.CreateUserResponse, error) {
    // insert user in users	
	var err error
	var response types.CreateUserResponse
	err = nil
	
	sql := fmt.Sprintf(`INSERT INTO users (user_id, username, 
	facebook_access_token, twitter_access_token, google_access_token, 
	push_token, device_type, email, password, is_test,
	is_anonymous, gender_preference, timestamp) 
	VALUES (%v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, %v,
	'%v', %v)`, 
	user.UserId, user.Username, user.FacebookAccessToken, user.TwitterAccessToken, 
	user.GoogleAccessToken, user.PushToken, user.DeviceType, 
	user.Email, user.Password, user.IsTest, user.IsAnonymous, user.GenderPreference, 
	user.Timestamp)
				
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraUserRepo.UpdateUser() - Error: %v", err.Error()))
	} else {
		response.UserId = user.UserId
		response.Username = user.Username
		response.Timestamp = user.Timestamp
	}

    return response, err
}

func (repo *CassandraUserRepo) GetUserByEmail(email string) (types.UserByEmail, error) {
	// get user by email	
	var user types.UserByEmail
	var user_id string
	var err error
	err = nil
	
	sql := fmt.Sprintf("SELECT email, user_id FROM users_by_email WHERE email = '%v' LIMIT 1", email)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&user.Email, &user_id); err != nil {
			log.Printf(fmt.Sprintf("CassandraUserRepo.GetUserByEmail() - Error: %v", err.Error()))
	} else {
		user.UserId = uuid.Parse(user_id)
	}
	
	return user, err
}

func (repo *CassandraUserRepo) GetUserExtraInfo(userId string) (types.UserExtraInfo, error) {
	// get user extra info	
	var user types.UserExtraInfo
	var user_id string
	var err error
	err = nil
	
	sql := fmt.Sprintf("SELECT user_id, walking_progress, timestamp FROM users_extra_info WHERE user_id = %v LIMIT 1", userId)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&user_id, &user.WalkingProgress, &user.Timestamp); err != nil {
		log.Printf(fmt.Sprintf("CassandraUserRepo.GetUserExtraInfo() - Error: %v", err.Error()))
	} else {
		user.UserId = uuid.Parse(user_id)
	}
	
	return user, err
}

func (repo *CassandraUserRepo) UpdateUserExtraInfo(ue types.UserExtraInfo) error {
	// update user extra info
	var err error
	err = nil
	
	sql := fmt.Sprintf(`INSERT INTO users_extra_info (user_id, walking_progress, timestamp) 
	values (%v, %v, %v)`, ue.UserId, ue.WalkingProgress, ue.Timestamp)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraUserRepo.UpdateUserExtraInfo() - Error: %v", err.Error()))
	}
	
	return err
}

func (repo *CassandraUserRepo) ReportUser(ur types.UserReport) error {
	// report user
	var err error
	err = nil
	
	sql := fmt.Sprintf(`INSERT INTO users_reported (reported_user_id, reporter_user_id, reason, comment, timestamp) 
	values (%v, %v, '%v', '%v', %v)`, ur.ReportedUserId, ur.ReporterUserId, ur.Reason, ur.Comment, ur.Timestamp)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraUserRepo.ReportUser() - Error: %v", err.Error()))
	}
	
	return err
}
