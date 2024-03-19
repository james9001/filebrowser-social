package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/notifications"
	"github.com/filebrowser/filebrowser/v2/usernotifications"
)

type notificationsResponse struct {
	Notifications []*notifications.Notification `json:"notifications"`
	MyUsername    string                        `json:"myUsername"`
}
type notificationsWithAckStateResponse struct {
	Notifications []*notificationWithAckState `json:"notifications"`
}

type notificationWithAckState struct {
	ID               uint      `storm:"id,increment" json:"id"`
	ContextFilePath  string    `json:"contextFilePath"`
	CausingUserName  string    `json:"causingUserName"`
	NotificationType string    `json:"notificationType"`
	CreatedTime      time.Time `json:"createdTime"`
	Acknowledged     bool      `json:"acknowledged"`
}

type unacknowledgedCountResponse struct {
	Count int
}

type acknowledgeNotificationRequest struct {
	modifyRequest
	Data *acknowledgeNotificationData `json:"data"`
}

type acknowledgeNotificationData struct {
	ContextFilePath string `json:"contextFilePath"`
}

type userUploadedNotificationRequest struct {
	modifyRequest
	Data *userUploadedNotificationData `json:"data"`
}

type userUploadedNotificationData struct {
	ContextFilePath string `json:"contextFilePath"`
}

func deserialiseAcknowledgeNotificationRequest(_ http.ResponseWriter, httpRequest *http.Request) (*acknowledgeNotificationRequest, error) {
	if httpRequest.Body == nil {
		return nil, errors.ErrEmptyRequest
	}

	req := &acknowledgeNotificationRequest{}
	err := json.NewDecoder(httpRequest.Body).Decode(req)
	if err != nil {
		return nil, err
	}
	if req.What != "acknowledge_notification" {
		return nil, errors.ErrInvalidDataType
	}

	return req, nil
}

var notificationsallGetHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	notifications, persistenceErr := inputData.store.Notifications.Gets()
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}
	notificationsResponse := notificationsResponse{
		MyUsername:    inputData.user.Username,
		Notifications: notifications,
	}

	return renderJSON(responseWriter, httpRequest, notificationsResponse)
})

var notificationUnacknowledgedCountGetHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	usernotifications, persistenceErr := inputData.store.UserNotifications.FindByUserName(inputData.user.Username)
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}

	notifications, persistenceErr := inputData.store.Notifications.Gets()
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}
	count := 0
	for _, notification := range notifications {
		if !usernotifications.AcknowledgedNotifications[notification.ID] {
			count++
		}
	}

	countResponse := unacknowledgedCountResponse{
		Count: count,
	}

	return renderJSON(responseWriter, httpRequest, countResponse)
})

var notificationAcknowledgePostHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	req, err := deserialiseAcknowledgeNotificationRequest(responseWriter, httpRequest)
	if err != nil {
		return http.StatusBadRequest, err
	}
	usernotifications, persistenceErr := inputData.store.UserNotifications.FindByUserName(inputData.user.Username)
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}

	notificationsWithSameFileContextPath, _ := inputData.store.Notifications.FindByContextFilePath(req.Data.ContextFilePath)

	for _, notification := range notificationsWithSameFileContextPath {
		usernotifications.AcknowledgedNotifications[notification.ID] = true
	}
	err = inputData.store.UserNotifications.Save(usernotifications)

	return errToStatus(err), err
})

func serialiseNotificationsWithAckStateResponse(notifications []*notifications.Notification, usernotifications *usernotifications.UserNotifications) *notificationsWithAckStateResponse {
	var ackedNotifications []*notificationWithAckState

	for _, a := range notifications {
		b := &notificationWithAckState{
			ID:               a.ID,
			ContextFilePath:  a.ContextFilePath,
			CausingUserName:  a.CausingUserName,
			NotificationType: a.NotificationType,
			CreatedTime:      a.CreatedTime,
			Acknowledged:     usernotifications.AcknowledgedNotifications[a.ID],
		}
		ackedNotifications = append(ackedNotifications, b)
	}

	notificationsWithAckStateResponse := notificationsWithAckStateResponse{
		Notifications: ackedNotifications,
	}
	return &notificationsWithAckStateResponse
}

var notificationsGetPageHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	vars := mux.Vars(httpRequest)
	pageNum, _ := strconv.Atoi(vars["pagenum"])

	notifications, persistenceErr := inputData.store.Notifications.FetchPage(pageNum)
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}

	usernotifications, persistenceErr := inputData.store.UserNotifications.FindByUserName(inputData.user.Username)

	return renderJSON(responseWriter, httpRequest, serialiseNotificationsWithAckStateResponse(notifications, usernotifications))
})

var notificationsGetRangeHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	lowId, _ := strconv.Atoi(httpRequest.URL.Query().Get("lowId"))
	highId, _ := strconv.Atoi(httpRequest.URL.Query().Get("highId"))

	notifications, persistenceErr := inputData.store.Notifications.GetRange(lowId, highId)
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}

	usernotifications, persistenceErr := inputData.store.UserNotifications.FindByUserName(inputData.user.Username)

	return renderJSON(responseWriter, httpRequest, serialiseNotificationsWithAckStateResponse(notifications, usernotifications))
})

func deserialiseUserUploadedNotificationRequest(_ http.ResponseWriter, httpRequest *http.Request) (*userUploadedNotificationRequest, error) {
	if httpRequest.Body == nil {
		return nil, errors.ErrEmptyRequest
	}

	req := &userUploadedNotificationRequest{}
	err := json.NewDecoder(httpRequest.Body).Decode(req)
	if err != nil {
		return nil, err
	}
	if req.What != "user_uploaded_notification" {
		return nil, errors.ErrInvalidDataType
	}

	return req, nil
}

var userUploadedNotificationPostHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	req, err := deserialiseUserUploadedNotificationRequest(responseWriter, httpRequest)
	if err != nil {
		return http.StatusBadRequest, err
	}

	createUserUploadedNotification(inputData, req)

	return errToStatus(err), err
})

func createUserUploadedNotification(inputData *data, req *userUploadedNotificationRequest) {

	newNotification := &notifications.Notification{
		ContextFilePath:  req.Data.ContextFilePath,
		CausingUserName:  inputData.user.Username,
		NotificationType: "UserUploaded",
		CreatedTime:      time.Now(),
	}

	inputData.store.Notifications.Save(newNotification)

	usernotifications, _ := inputData.store.UserNotifications.FindByUserName(newNotification.CausingUserName)

	usernotifications.AcknowledgedNotifications[newNotification.ID] = true
	_ = inputData.store.UserNotifications.Save(usernotifications)

}
