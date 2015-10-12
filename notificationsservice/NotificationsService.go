// NotificationsService
package notificationsservice

import (
	"github.com/somanole/gaitapi/types"
	"errors"
	"strings"
)

func SendPushNotification(deviceType string, pushToken string, alertText string) error {
	var pn types.PushNotification
	var err error
	err = errors.New("400")
	
	pn.PushToken = pushToken
	pn.AlertText = alertText
	
	switch strings.ToLower(deviceType) {
		case "ios" :
			err = SendApplePushNotification(pn)
		case "android" :
			err = SendAndroidPushNotification(pn)
	}
	
	return err
}
