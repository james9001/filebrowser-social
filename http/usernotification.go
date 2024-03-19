package http

import (
	"net/http"
)

var usernotificationsGetAllHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	usernotifications, persistenceErr := inputData.store.UserNotifications.Gets()
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}

	return renderJSON(responseWriter, httpRequest, usernotifications)
})

var usernotificationsGetAllMineHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	usernotifications, persistenceErr := inputData.store.UserNotifications.FindByUserName(inputData.user.Username)
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}

	return renderJSON(responseWriter, httpRequest, usernotifications)
})
