// CassandraMessageRepo
package messagerepo

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
	CassandraMessageRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraMessageRepo() MessageRepo {
	return &CassandraMessageRepo{}
}

func init() {
	New = NewCassandraMessageRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraMessageRepo) CreateMessage(userId string, receiverId string, mr types.MessageRequest) error {
    // insert message in messages
	var m types.Message
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil && uuid.Parse(receiverId) != nil {
		sql := fmt.Sprintf(`SELECT email from users_by_id WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		var email string
		if err = session.Query(sql).Scan(&email); err != nil {
				log.Printf(fmt.Sprintf("CreateMessage - Error: %v", err.Error()))
		} else {
			m.SenderId = uuid.Parse(userId)
			m.ReceiverId = uuid.Parse(receiverId)
			m.MessageId = uuid.NewRandom()
			m.IsRead = false
			m.Text = mr.Text
			if mr.Timestamp != 0 {
				m.Timestamp = int64(time.Unix(mr.Timestamp, 0).UTC().Unix())
			} else {
				m.Timestamp = int64(time.Now().UTC().Unix())
			}
			
			sql := fmt.Sprintf(`INSERT INTO messages (message_id, sender_id, 
			receiver_id, text, is_read, timestamp) VALUES (%v, %v, %v, '%v', %v, %v)`, 
			m.MessageId, m.SenderId, m.ReceiverId, m.Text, m.IsRead, m.Timestamp)
						
			log.Printf(sql)
			if err = session.Query(sql).Exec(); err != nil {
				log.Printf(fmt.Sprintf("CreateMessage - Error: %v", err.Error()))
			} 
		}
	} else {
		log.Printf(fmt.Sprintf("GetUserMatch - Sender | Received Id: %v | %v is not UUID", userId, m.ReceiverId))
		err = errors.New("not uuid")
	}
	
    return err
}

func (repo *CassandraMessageRepo) GetUserMessagesByReceiverId(userId string, receiverId string, startdate string) (types.Messages, error) {
	// get all messages for user
	log.Printf(fmt.Sprintf("GetUserMessages - Received userId: %v", userId))
	log.Printf(fmt.Sprintf("GetUserMessages - Received receiverId: %v", receiverId))
	
	var messages types.Messages
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil && uuid.Parse(receiverId) != nil {
		var message_id, sender_id, receiver_id string
		var text string
		var is_read bool
		var timestamp int64
		var sql string
		
		if startdate != "" {
			sql = fmt.Sprintf(`SELECT message_id, sender_id, receiver_id, text, is_read, 
			timestamp FROM messages WHERE sender_id = %v AND receiver_id = %v AND timestamp >= %v`, userId, receiverId, startdate)
		} else {
			sql = fmt.Sprintf(`SELECT message_id, sender_id, receiver_id, text, is_read, 
			timestamp FROM messages WHERE sender_id = %v AND receiver_id = %v`, userId, receiverId)
		}
		
		log.Printf(sql)
		
		iter := session.Query(sql).Iter()
		for iter.Scan(&message_id, &sender_id, &receiver_id, &text, &is_read, &timestamp) {
			messages = append(messages, types.Message{MessageId: uuid.Parse(message_id), SenderId: uuid.Parse(sender_id), 
			ReceiverId: uuid.Parse(receiver_id), Text: text, IsRead: is_read, Timestamp: timestamp})
		}
		if err = iter.Close(); err != nil {
			log.Printf(fmt.Sprintf("GetUserMessages - Error: %v", err.Error()))
		}
	} else {
		log.Printf(fmt.Sprintf("GetUserMatch - Received userId | receiverId : %v | %v is not UUID", userId, receiverId))
		err = errors.New("not uuid")
	}
	
	return messages, err
}