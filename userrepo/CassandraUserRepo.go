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

func (repo *CassandraUserRepo) CreateUser(u types.User) types.CreateUserResponse {
    // insert user in users
	u.UserId = uuid.NewRandom()
	u.Username = "brown fox drinks wine"
	
	u.Timestamp = int64(time.Now().UTC().Unix())

	sql :=fmt.Sprintf(`INSERT INTO users (user_id, username, 
	facebook_access_token, twitter_access_token, google_access_token, 
	push_token_ios, push_token_android, device_type, email, password, is_test,
	is_anonymous, gender_preference, timestamp) 
	VALUES (%v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, %v,
	'%v', %v)`, 
	u.UserId, u.Username, u.FacebookAccessToken, u.TwitterAccessToken, 
	u.GoogleAccessToken, u.PushTokeniOS, u.PushTokenAndroid, u.DeviceType, 
	u.Email, u.Password, u.IsTest, u.IsAnonymous, u.GenderPreference, 
	u.Timestamp)
	
	log.Printf(sql)
	if err := session.Query(sql).Exec(); err != nil {
		log.Fatal(err)
	}
	
	var response types.CreateUserResponse
	response.UserId = u.UserId
	response.Username = u.Username
	
    return response
}

func (repo *CassandraUserRepo) GetUser(userId string) (types.User, error) {
	// get user by id
	log.Printf(fmt.Sprintf("Received userId: %v", userId))
	
	var user types.User
	var user_id string
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil{
		sql := fmt.Sprintf(`SELECT user_id, username, facebook_access_token, 
		twitter_access_token, google_access_token, push_token_ios, push_token_android, 
		device_type, email, password, is_test, is_anonymous, gender_preference, 
		timestamp FROM users WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&user_id, 
			&user.Username, &user.FacebookAccessToken, &user.TwitterAccessToken, &user.GoogleAccessToken,
			&user.PushTokeniOS, &user.PushTokenAndroid, &user.DeviceType, &user.Email,
			&user.Password, &user.IsTest, &user.IsAnonymous, &user.GenderPreference, &user.Timestamp); err != nil {
				log.Printf(err.Error())
		}
		
		user.UserId = uuid.Parse(user_id)
	} else{
		log.Printf(fmt.Sprintf("Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return user, err
}

func (repo *CassandraUserRepo) UpdateUser(userId string, u types.User) (types.CreateUserResponse, error) {
    // insert user in users
	log.Printf(fmt.Sprintf("Received userId: %v", userId))
	
	var user types.User
	var err error
	var response types.CreateUserResponse
	err = nil
	
	if uuid.Parse(userId) != nil {
		var user_id string
		
		sql := fmt.Sprintf(`SELECT user_id, username, facebook_access_token, 
		twitter_access_token, google_access_token, push_token_ios, push_token_android, 
		device_type, email, password, is_test, is_anonymous, gender_preference, 
		timestamp FROM users WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&user_id, 
			&user.Username, &user.FacebookAccessToken, &user.TwitterAccessToken, &user.GoogleAccessToken,
			&user.PushTokeniOS, &user.PushTokenAndroid, &user.DeviceType, &user.Email,
			&user.Password, &user.IsTest, &user.IsAnonymous, &user.GenderPreference, &user.Timestamp); err != nil {
				log.Printf(err.Error())
		} else {
			user.UserId = uuid.Parse(user_id)
			user.Timestamp = int64(time.Now().UTC().Unix())
			
			if u.FacebookAccessToken != "" { 
				user.FacebookAccessToken = u.FacebookAccessToken
			}
			if u.DeviceType != "" {
				user.DeviceType = u.DeviceType
			}
			if u.Email != "" {
				user.Email = u.Email
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
			if u.PushTokenAndroid != "" {
				user.PushTokenAndroid = u.PushTokenAndroid
			}
			if u.PushTokeniOS != "" {
				user.PushTokeniOS = u.PushTokenAndroid
			}
			if u.TwitterAccessToken != "" {
				user.TwitterAccessToken = u.TwitterAccessToken
			}
			
			sql := fmt.Sprintf(`INSERT INTO users (user_id, username, 
			facebook_access_token, twitter_access_token, google_access_token, 
			push_token_ios, push_token_android, device_type, email, password, is_test,
			is_anonymous, gender_preference, timestamp) 
			VALUES (%v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, %v,
			'%v', %v)`, 
			user.UserId, user.Username, user.FacebookAccessToken, user.TwitterAccessToken, 
			user.GoogleAccessToken, user.PushTokeniOS, user.PushTokenAndroid, user.DeviceType, 
			user.Email, user.Password, user.IsTest, user.IsAnonymous, user.GenderPreference, 
			user.Timestamp)
				
			log.Printf(sql)
			if err = session.Query(sql).Exec(); err != nil {
				log.Printf(err.Error())
			}
				
			response.UserId = user.UserId
			response.Username = user.Username
		}
	} else {
		log.Printf(fmt.Sprintf("Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}

    return response, err
}
