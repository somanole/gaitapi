// CassandraUtilsRepo
package utilsrepo

import (
	"errors"
	"fmt"
	"log"
	"github.com/gocql/gocql"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/constants"
)

type (
	CassandraUtilsRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraUtilsRepo() UtilsRepo {
	return &CassandraUtilsRepo{}
}

func init() {
	New = NewCassandraUtilsRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster(constants.CASSANDRA_MASTER)
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraUtilsRepo) CheckUserPassword(userId string, rpassword string) error {
	// check password for user
	var password string
	var err error
	err = nil
	
	sql := fmt.Sprintf("SELECT password FROM users WHERE user_id = %v LIMIT 1", userId)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&password); err != nil {
			log.Printf(fmt.Sprintf("CassandraUtilsRepo.CheckUserPassword() - Error: %v", err.Error()))
	} else if password != rpassword{
		log.Printf(fmt.Sprintf("CassandraUtilsRepo.CheckUserPassword() - password received: %v does not match password on record: %v for userId: %v", rpassword, password, userId))
		err = errors.New("401")
	}
	
	return err
}

func (repo *CassandraUtilsRepo) CheckIfUserExists(userId string) error {
	// check if user exists
	var password string
	var err error
	err = nil
	
	sql := fmt.Sprintf("SELECT password FROM users WHERE user_id = %v LIMIT 1", userId)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&password); err != nil {
		log.Printf(fmt.Sprintf("CassandraUtilsRepo.CheckIfUserExists() - Error: %v", err.Error()))
	}

	return err
}

func (repo *CassandraUtilsRepo) CheckIfMatchExists(firstUserId string, secondUserId string) error {
	// check if user exists
	var matched_username string
	var err error
	err = nil
	
	sql := fmt.Sprintf(`SELECT matched_username FROM matches_by_matched_user_id 
	WHERE user_id = %v and matched_user_id = %v LIMIT 1`, firstUserId, secondUserId)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&matched_username); err != nil {
		log.Printf(fmt.Sprintf("CassandraUtilsRepo.CheckIfMatchExists() - Error: %v", err.Error()))
	} else {
		sql = fmt.Sprintf(`SELECT matched_username FROM matches_by_matched_user_id 
		WHERE user_id = %v and matched_user_id = %v LIMIT 1`, secondUserId, firstUserId)
			
		log.Printf(sql)
			
		if err = session.Query(sql).Scan(&matched_username); err != nil {
			log.Printf(fmt.Sprintf("CassandraUtilsRepo.CheckIfMatchExists() - Error: %v", err.Error()))
		}
	}

	return err
}

func (repo *CassandraUtilsRepo) GetUserUsername(userId string) (string, error) {
	// get user username
	var username string
	var err error
	err = nil
	
	sql := fmt.Sprintf("SELECT username FROM users WHERE user_id = %v LIMIT 1", userId)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&username); err != nil {
		log.Printf(fmt.Sprintf("CassandraUtilsRepo.GetUserUsername() - Error: %v", err.Error()))
	}
	
	return username, err
}

func (repo *CassandraUtilsRepo) CheckLoginCredentials(email string, rpassword string) (string, error) {
	// check login credentials
	var userId string
	var err error
	err = nil
	
	sql := fmt.Sprintf("SELECT user_id FROM users_by_email WHERE email = '%v' LIMIT 1", email)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&userId); err != nil {
		log.Printf(fmt.Sprintf("CassandraUtilsRepo.CheckLoginCredentials() - Error: %v", err.Error()))
	} else {
		var password string
		
		sql = fmt.Sprintf("SELECT password FROM users WHERE user_id = %v LIMIT 1", userId)
		
		log.Printf(sql)
			
		if err = session.Query(sql).Scan(&password); err != nil {
			log.Printf(fmt.Sprintf("CassandraUtilsRepo.CheckLoginCredentials() - Error: %v", err.Error()))
		} else if password != rpassword{
			log.Printf(fmt.Sprintf("CassandraUtilsRepo.CheckLoginCredentials() - password received: %v does not match password on record: %v for email: %v", rpassword, password, email))
			err = errors.New("401")
		}
	}
	
	return userId, err
}

func (repo *CassandraUtilsRepo) RegisterInterest(i types.Interest) error {
    // insert interest
	var err error
	err = nil
	
	sql :=fmt.Sprintf("INSERT INTO interests (email, timestamp) VALUES ('%v', %v)", i.Email, i.Timestamp)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraUtilsRepo.RegisterInterest() - Error: %v", err.Error()))
	} 
	
    return err
}
