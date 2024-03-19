package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/notifications"
	"github.com/filebrowser/filebrowser/v2/reactions"
)

type reactionsResponse struct {
	Reactions  []*reactions.Reaction `json:"reactions"`
	MyUsername string                `json:"myUsername"`
}

type reactionAvailableListResponse struct {
	ReactionsAvailable []string `json:"reactionsAvailable"`
}

type createReactionRequest struct {
	modifyRequest
	Data *reactions.Reaction `json:"data"`
}

var reactionGetAvailableReactionList = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	file, err := files.NewFileInfo(files.FileOptions{
		Fs:         inputData.user.Fs,
		Path:       "/filebrowser-social/reactions/",
		Modify:     inputData.user.Perm.Modify,
		Expand:     true,
		ReadHeader: inputData.server.TypeDetectionByHeader,
		Checker:    inputData,
		Content:    true,
	})

	if !file.IsDir {
		return errToStatus(err), err
	}

	var availableReactions []string

	for _, a := range file.Items {
		reactionIndex := strings.Index(a.Path, ".png")

		reaction := a.Path[len("/filebrowser-social/reactions/"):reactionIndex]

		availableReactions = append(availableReactions, reaction)
	}

	reactionAvailableListResponse := reactionAvailableListResponse{
		ReactionsAvailable: availableReactions,
	}
	return renderJSON(responseWriter, httpRequest, &reactionAvailableListResponse)

})

func deserialiseCreateReactionRequest(_ http.ResponseWriter, httpRequest *http.Request) (*createReactionRequest, error) {
	if httpRequest.Body == nil {
		return nil, errors.ErrEmptyRequest
	}

	req := &createReactionRequest{}
	err := json.NewDecoder(httpRequest.Body).Decode(req)
	if err != nil {
		return nil, err
	}
	if req.What != "create_reaction" {
		return nil, errors.ErrInvalidDataType
	}

	return req, nil
}

func createReactionNotification(inputData *data, reaction *reactions.Reaction, reactionType string) {

	newNotification := &notifications.Notification{
		ContextFilePath:  reaction.ContextFilePath,
		CausingUserName:  reaction.UserName,
		NotificationType: reactionType,
		CreatedTime:      time.Now(),
	}

	inputData.store.Notifications.Save(newNotification)

	usernotifications, _ := inputData.store.UserNotifications.FindByUserName(reaction.UserName)

	usernotifications.AcknowledgedNotifications[newNotification.ID] = true
	_ = inputData.store.UserNotifications.Save(usernotifications)

}

var reactionPostHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	req, err := deserialiseCreateReactionRequest(responseWriter, httpRequest)
	if err != nil {
		return http.StatusBadRequest, err
	}

	req.Data.UserName = inputData.user.Username
	req.Data.CreatedTime = time.Now()

	err = inputData.store.Reactions.Save(req.Data)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	notificationType := func() string {
		if req.Data.CommentID == 0 {
			return "UserReactsToFile"
		}
		return "UserReactsToComment"
	}()
	createReactionNotification(inputData, req.Data, notificationType)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return errToStatus(err), err
})

var reactionsDebugGetAllHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	reactions, persistenceErr := inputData.store.Reactions.Gets()
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}
	reactionsResponse := reactionsResponse{
		MyUsername: inputData.user.Username,
		Reactions:  reactions,
	}

	return renderJSON(responseWriter, httpRequest, reactionsResponse)
})

var reactionsGetByContextFilePathHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {

	reactions, persistenceErr := inputData.store.Reactions.FindByContextFilePath(httpRequest.URL.Path)
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}
	reactionsResponse := reactionsResponse{
		MyUsername: inputData.user.Username,
		Reactions:  reactions,
	}

	return renderJSON(responseWriter, httpRequest, reactionsResponse)
})

var reactionDeleteHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {
	reactionId := strings.TrimSuffix(httpRequest.URL.Path, "/")
	reactionId = strings.TrimPrefix(reactionId, "/")

	reactionIdUint64, _ := strconv.ParseUint(reactionId, 10, 32)
	reactionIdUint := uint(reactionIdUint64)

	reaction, _ := inputData.store.Reactions.GetById(reactionIdUint)
	if reaction.UserName != inputData.user.Username {
		return errToStatus(errors.ErrPermissionDenied), errors.ErrPermissionDenied
	}

	err := inputData.store.Reactions.Delete(reactionIdUint)
	if err != nil {
		return errToStatus(err), err
	}

	return http.StatusOK, nil
})
