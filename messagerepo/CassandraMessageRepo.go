// CassandraMessageRepo
package messagerepo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/constants"
	"github.com/gocql/gocql"
	"code.google.com/p/go-uuid/uuid"
	"sort"
)

// ByTimestamp implements sort.Interface for []Message based on
// the Timestamp field.
type ByTimestamp []types.Message

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }

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
	cluster := gocql.NewCluster(constants.CASSANDRA_MASTER)
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraMessageRepo) CreateMessage(m types.Message) error {
    // insert message in messages
	var err error
	err = nil
	
	sql := fmt.Sprintf(`INSERT INTO messages (message_id, sender_id, 
	receiver_id, text, is_read, timestamp) VALUES (%v, %v, %v, '%v', %v, %v)`, 
	m.MessageId, m.SenderId, m.ReceiverId, m.Text, m.IsRead, m.Timestamp)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraMessageRepo.CreateMessage() - Error: %v", err.Error()))
	} 

    return err
}

func (repo *CassandraMessageRepo) GetUserMessagesByReceiverId(userId string, receiverId string, startdate string) (types.Messages, error) {
	// get all messages for user
	log.Printf(fmt.Sprintf("CassandraMessageRepo.GetUserMessages() - Received userId: %v", userId))
	log.Printf(fmt.Sprintf("CassandraMessageRepo.GetUserMessages() - Received receiverId: %v", receiverId))
	
	var messages types.Messages
	var err error
	err = nil
	
	var message_id, sender_id, receiver_id string
	var text string
	var is_read bool
	var timestamp int64
	var sql string
		
	//sender messages
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
		log.Printf(fmt.Sprintf("CassandraMessageRepo.GetUserMessages() - Error: %v", err.Error()))
	}
		
	//receiver messages
	if startdate != "" {
		sql = fmt.Sprintf(`SELECT message_id, sender_id, receiver_id, text, is_read, 
		timestamp FROM messages WHERE sender_id = %v AND receiver_id = %v AND timestamp >= %v`, receiverId, userId, startdate)
	} else {
		sql = fmt.Sprintf(`SELECT message_id, sender_id, receiver_id, text, is_read, 
		timestamp FROM messages WHERE sender_id = %v AND receiver_id = %v`, receiverId, userId)
	}
		
	log.Printf(sql)
		
	iter2 := session.Query(sql).Iter()
	for iter2.Scan(&message_id, &sender_id, &receiver_id, &text, &is_read, &timestamp) {
		messages = append(messages, types.Message{MessageId: uuid.Parse(message_id), SenderId: uuid.Parse(sender_id), 
		ReceiverId: uuid.Parse(receiver_id), Text: text, IsRead: is_read, Timestamp: timestamp})
	}
	if err = iter2.Close(); err != nil {
		log.Printf(fmt.Sprintf("CassandraMessageRepo.GetUserMessages() - Error: %v", err.Error()))
	}
		
	if messages != nil {
		sort.Sort(ByTimestamp(messages))
	}
	
	return messages, err
}

func (repo *CassandraMessageRepo) GetUserLastMessageByReceiverId(userId string, receiverId string) (types.Message, error) {
	// get last message for user
	log.Printf(fmt.Sprintf("CassandraMessageRepo.GetUserLastMessageByReceiverId() - Received userId: %v", userId))
	log.Printf(fmt.Sprintf("CassandraMessageRepo.GetUserLastMessageByReceiverId() - Received receiverId: %v", receiverId))
	
	var messages types.Messages
	var err error
	err = nil
	
	var message_id, sender_id, receiver_id string
	var text string
	var is_read bool
	var timestamp int64
	var sql string
		
	//sender messages
	sql = fmt.Sprintf(`SELECT message_id, sender_id, receiver_id, text, is_read, 
	timestamp FROM messages WHERE sender_id = %v AND receiver_id = %v LIMIT 1`, userId, receiverId)
		
	log.Printf(sql)
		
	iter := session.Query(sql).Iter()
	for iter.Scan(&message_id, &sender_id, &receiver_id, &text, &is_read, &timestamp) {
		messages = append(messages, types.Message{MessageId: uuid.Parse(message_id), SenderId: uuid.Parse(sender_id), 
		ReceiverId: uuid.Parse(receiver_id), Text: text, IsRead: is_read, Timestamp: timestamp})
	}
	if err = iter.Close(); err != nil {
		log.Printf(fmt.Sprintf("CassandraMessageRepo.GetUserLastMessageByReceiverId() - Error: %v", err.Error()))
	}
		
	//receiver messages
	sql = fmt.Sprintf(`SELECT message_id, sender_id, receiver_id, text, is_read, 
	timestamp FROM messages WHERE sender_id = %v AND receiver_id = %v`, receiverId, userId)
	
	log.Printf(sql)
		
	iter2 := session.Query(sql).Iter()
	for iter2.Scan(&message_id, &sender_id, &receiver_id, &text, &is_read, &timestamp) {
		messages = append(messages, types.Message{MessageId: uuid.Parse(message_id), SenderId: uuid.Parse(sender_id), 
		ReceiverId: uuid.Parse(receiver_id), Text: text, IsRead: is_read, Timestamp: timestamp})
	}
	if err = iter2.Close(); err != nil {
		log.Printf(fmt.Sprintf("CassandraMessageRepo.GetUserLastMessageByReceiverId() - Error: %v", err.Error()))
	}
		
	if messages != nil {
		sort.Sort(ByTimestamp(messages))
	}
	
	var message types.Message
	
	if len(messages) > 0 {
		message = messages[0]
	}
	
	return message, err
}

func (repo *CassandraMessageRepo) DeleteMessage(senderId uuid.UUID, receiverId uuid.UUID, timestamp int64) error {
	// get all messages for user
	log.Printf(fmt.Sprintf("CassandraMessageRepo.DeleteMessage() - Received senderId: %v", senderId))
	log.Printf(fmt.Sprintf("CassandraMessageRepo.DeleteMessage() - Received receiverId: %v", receiverId))
	log.Printf(fmt.Sprintf("CassandraMessageRepo.DeleteMessage() - Received timestamp: %v", timestamp))
	
	var err error
	err = nil
	
	sql := fmt.Sprintf("DELETE FROM messages WHERE sender_id = %v AND receiver_id = %v AND timestamp = %v", senderId, receiverId, timestamp)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraMessageRepo.DeleteMessage() - Error: %v", err.Error()))
	} 

    return err
}