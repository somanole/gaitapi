// +build !test

// CassandraUserRepo
package userrepo

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
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraUserRepo) CreateUser(ur types.UserRequest) (types.CreateUserResponse, error) {
    // insert user in users
	var u types.User
	var response types.CreateUserResponse
	var existingEmail string
	var err error
	var sql string
	err = nil
	existingEmail = ""
	
	if (!ur.IsAnonymous) {
		sql = fmt.Sprintf("SELECT email FROM users_by_email WHERE email = '%v' LIMIT 1", ur.Email)		
		log.Printf(sql)
	
		if err = session.Query(sql).Scan(&existingEmail); err != nil {
			log.Printf(fmt.Sprintf("CreateUser - Error: %v", err.Error()))
		}
	}
	
	if (ur.IsAnonymous || (!ur.IsAnonymous && err != nil && err.Error() == "not found")) {
		err = nil
				
		u.UserId = uuid.NewRandom()
		u.Username = "brown fox drinks wine"	
		u.Timestamp = int64(time.Now().UTC().Unix())
		u.DeviceType = ur.DeviceType
		u.Email = ur.Email
		u.FacebookAccessToken = ur.FacebookAccessToken
		u.GenderPreference = ur.GenderPreference
		u.GoogleAccessToken = ur.GoogleAccessToken
		u.IsAnonymous = ur.IsAnonymous
		u.IsTest = false
		u.Password = ur.Password
		u.PushToken = ur.PushToken
		u.TwitterAccessToken = ur.TwitterAccessToken
				
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
			log.Printf(fmt.Sprintf("CreateUser - Error: %v", err.Error()))
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
				log.Printf(fmt.Sprintf("CreateUser - Error: %v", err.Error()))
			} 
		}
	} else if existingEmail != "" {
		err = errors.New("email already registered")
	} 
	
    return response, err
}

func (repo *CassandraUserRepo) GetUser(userId string) (types.User, error) {
	// get user by id
	log.Printf(fmt.Sprintf("GetUser - Received userId: %v", userId))
	
	var user types.User
	var user_id string
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf(`SELECT user_id, username, facebook_access_token, 
		twitter_access_token, google_access_token, push_token, 
		device_type, email, password, is_test, is_anonymous, gender_preference, 
		timestamp FROM users WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&user_id, 
			&user.Username, &user.FacebookAccessToken, &user.TwitterAccessToken, &user.GoogleAccessToken,
			&user.PushToken, &user.DeviceType, &user.Email,
			&user.Password, &user.IsTest, &user.IsAnonymous, &user.GenderPreference, &user.Timestamp); err != nil {
				log.Printf(fmt.Sprintf("GetUser - Error: %v", err.Error()))
		} else {
			user.UserId = uuid.Parse(user_id)
		}
	} else {
		log.Printf(fmt.Sprintf("GetUser - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return user, err
}

func (repo *CassandraUserRepo) UpdateUser(userId string, u types.UserUpdateRequest) (types.CreateUserResponse, error) {
    // insert user in users
	log.Printf(fmt.Sprintf("UpdateUser - Received userId: %v", userId))
	
	var user types.User
	var err error
	var response types.CreateUserResponse
	err = nil
	
	if uuid.Parse(userId) != nil {
		var user_id string
		
		sql := fmt.Sprintf(`SELECT user_id, username, facebook_access_token, 
		twitter_access_token, google_access_token, push_token, 
		device_type, email, password, is_test, is_anonymous, gender_preference, 
		timestamp FROM users WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&user_id, 
			&user.Username, &user.FacebookAccessToken, &user.TwitterAccessToken, &user.GoogleAccessToken,
			&user.PushToken, &user.DeviceType, &user.Email,
			&user.Password, &user.IsTest, &user.IsAnonymous, &user.GenderPreference, &user.Timestamp); err != nil {
				log.Printf(fmt.Sprintf("UpdateUser - Error: %v", err.Error()))
		} else {
			user.UserId = uuid.Parse(user_id)
			user.Timestamp = int64(time.Now().UTC().Unix())
			
			if u.FacebookAccessToken != "" { 
				user.FacebookAccessToken = u.FacebookAccessToken
			}
			if u.DeviceType != "" {
				user.DeviceType = u.DeviceType
			}
			if u.GenderPreference != "" {
				user.GenderPreference = u.GenderPreference
			}
			if u.GoogleAccessToken != "" {
				user.GoogleAccessToken = u.GoogleAccessToken
			}
			if u.Password != "" {
				user.Password = u.Password
			}
			if u.PushToken != "" {
				user.PushToken = u.PushToken
			}
			if u.TwitterAccessToken != "" {
				user.TwitterAccessToken = u.TwitterAccessToken
			}
			
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
				log.Printf(fmt.Sprintf("UpdateUser - Error: %v", err.Error()))
			} else {
				response.UserId = user.UserId
				response.Username = user.Username
				response.Timestamp = user.Timestamp
			}	
		}
	} else {
		log.Printf(fmt.Sprintf("UpdateUser - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}

    return response, err
}

func (repo *CassandraUserRepo) GetUserByEmail(email string) (types.UserByEmail, error) {
	// get user by email
	log.Printf(fmt.Sprintf("GetUserByEmail - Received email: %v", email))
	
	var user types.UserByEmail
	var user_id string
	var err error
	err = nil
	
	sql := fmt.Sprintf("SELECT email, user_id FROM users_by_email WHERE email = '%v' LIMIT 1", email)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&user.Email, &user_id); err != nil {
			log.Printf(fmt.Sprintf("GetUserByEmail - Error: %v", err.Error()))
	} else {
		user.UserId = uuid.Parse(user_id)
	}
	
	return user, err
}

func (repo *CassandraUserRepo) GetUserExtraInfo(userId string) (types.UserExtraInfo, error) {
	// get user extra info
	log.Printf(fmt.Sprintf("GetUserExtraInfo - Received userId: %v", userId))
	
	var user types.UserExtraInfo
	var user_id string
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf("SELECT user_id, walking_progress, timestamp FROM users_extra_info WHERE user_id = %v LIMIT 1", userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&user_id, &user.WalkingProgress, &user.Timestamp); err != nil {
				log.Printf(fmt.Sprintf("GetUserExtraInfo - Error: %v", err.Error()))
		} else {
			user.UserId = uuid.Parse(user_id)
		}
	} else {
		log.Printf(fmt.Sprintf("GetUserExtraInfo - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return user, err
}
