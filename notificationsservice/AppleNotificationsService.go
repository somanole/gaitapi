// AppleNotificationsService
package notificationsservice

import (
  	apns "github.com/joekarl/go-libapns"
  	"fmt"
  	"io/ioutil"
  	"time"
	"log"
	"github.com/somanole/gaitapi/types"
)

var config *apns.APNSConfig = getAPNSConfig()

func getAPNSConfig() *apns.APNSConfig {
	certPem, err := ioutil.ReadFile("certificates/apns/cert.pem")
	if err != nil {
	   	log.Printf(fmt.Sprintf("AppleNotificationsService.getAPNSConnection() - Error: %v", err.Error()))
	}
	keyPem, err := ioutil.ReadFile("certificates/apns/key.pem")
	if err != nil {
	   	log.Printf(fmt.Sprintf("AppleNotificationsService.getAPNSConnection() - Error: %v", err.Error()))
	}
	
	config := &apns.APNSConfig{
	  	CertificateBytes: certPem,
	    KeyBytes: keyPem,
	    GatewayHost: "gateway.push.apple.com",
	}
	
	return config
}

type HandleCloseErrorFn func (closeError *apns.ConnectionClose)

func poster(sendChannel chan<- *apns.Payload, doneWritingChan chan bool, pushToken string, alertText string) {
    payload := &apns.Payload {
      Token: pushToken,
      AlertText: alertText,
    }
	
    sendChannel <- payload

  	doneWritingChan <- true
}

func apnsService(config *apns.APNSConfig, sendChannel <-chan *apns.Payload, closeErrFn HandleCloseErrorFn, shutdownChannel chan bool) {
	shutdown := false
	connectionGood := false
	var lastConn *apns.APNSConnection = nil
	
	for {
	    if shutdown {
	      break
	    }
		
	    log.Printf("establishing connection")
	    conn, err := apns.NewAPNSConnection(config)
		
	    if err != nil {
	       connectionGood = false
	       fmt.Println(err)
		
	      	select {
	        	case <-time.After(time.Second * 5):
	          		continue
	        	case <-shutdownChannel:
	          		shutdown = true
	      	}
	    } else {
	      	lastConn = conn
	      	connectionGood = true
	    }
	
	    for {
	      	if !connectionGood || shutdown {
	        	break
	      	}
			
	      	select {
	        	case payload := <- sendChannel:
	          		log.Printf("sending id %v\n", payload.ExtraData)
	          		select {
	            		case <-time.After(time.Second * 1):
	              			break
	            		case conn.SendChannel <- payload:
	              			break
	          		}
	          		break
	        	case closeError := <-conn.CloseChannel:
	          		closeErrFn(closeError)
	          		connectionGood = false
	          		log.Printf("Received error, closing connection")
	          		break
	        	case <-shutdownChannel:
	          		log.Printf("Received shutdown signal, closing connection")
	          		conn.Disconnect()
	          		shutdown = true
	      	}
	    }
	    log.Printf("Connection killed attempting re-establish")
	}
	
	if lastConn != nil {
	    select {
	    	case <-time.After(time.Second * 5):
	        	break
	      	case closeError := <-lastConn.CloseChannel:
	        	closeErrFn(closeError)
	    }
	}
	
	log.Printf("going to shutdown")
	shutdownChannel <- true
}

func handleCloseError(closeError *apns.ConnectionClose) {
  	// do something here with unsent payloads
  	log.Printf(fmt.Sprintf("AppleNotificationsService.handleCloseError() Error: +%v", closeError.Error))
}



func SendApplePushNotification(pn types.PushNotification) error {
	log.Printf(fmt.Sprintf("AppleNotificationsService.SendApplePushNotification() - PushToken: %v, AlertText: %v", pn.PushToken, pn.AlertText))
	var err error
	err = nil
	
	shutdownChannel := make(chan bool)
	sendChannel := make(chan *apns.Payload)
	doneWritingChan := make(chan bool)
	
	go poster(sendChannel, doneWritingChan, pn.PushToken, pn.AlertText)
	go apnsService(config, sendChannel, handleCloseError, shutdownChannel)
	
	// wait till we're done writing
	<-doneWritingChan
	shutdownChannel <- true
	
	// wait for things to shutdown
	log.Printf("waiting to shutdown")
	<-shutdownChannel
	
	return err
}
