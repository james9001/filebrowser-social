package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/filebrowser/filebrowser/v2/comments"
	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/notifications"
)

type modifyCommentRequest struct {
	modifyRequest
	Data *comments.Comment `json:"data"`
}

type commentsResponse struct {
	Comments   []*comments.Comment `json:"comments"`
	MyUsername string              `json:"myUsername"`
}

func deserialiseComment(_ http.ResponseWriter, httpRequest *http.Request) (*modifyCommentRequest, error) {
	if httpRequest.Body == nil {
		return nil, errors.ErrEmptyRequest
	}

	req := &modifyCommentRequest{}
	err := json.NewDecoder(httpRequest.Body).Decode(req)
	if err != nil {
		return nil, err
	}
	if req.What != "comment" {
		return nil, errors.ErrInvalidDataType
	}

	return req, nil
}

var commentGetForFileHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {
	//TODO: this is here to validate user has access to file, should make this consistent
	_, fileErr := files.NewFileInfo(files.FileOptions{
		Fs:         inputData.user.Fs,
		Path:       httpRequest.URL.Path,
		Modify:     inputData.user.Perm.Modify,
		Expand:     true,
		ReadHeader: inputData.server.TypeDetectionByHeader,
		Checker:    inputData,
		Content:    true,
	})
	if fileErr != nil {
		return errToStatus(fileErr), fileErr
	}

	filePath := httpRequest.URL.Path
	comments, persistenceErr := inputData.store.Comments.FindByFilePath(filePath)
	if persistenceErr != nil {
		return errToStatus(persistenceErr), persistenceErr
	}
	commentsResponse := commentsResponse{
		MyUsername: inputData.user.Username,
		Comments:   comments,
	}

	return renderJSON(responseWriter, httpRequest, commentsResponse)
})

var commentPostHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {
	//TODO: this is here to validate user has access to file, should make this consistent
	_, fileErr := files.NewFileInfo(files.FileOptions{
		Fs:         inputData.user.Fs,
		Path:       httpRequest.URL.Path,
		Modify:     inputData.user.Perm.Modify,
		Expand:     true,
		ReadHeader: inputData.server.TypeDetectionByHeader,
		Checker:    inputData,
		Content:    true,
	})
	if fileErr != nil {
		return errToStatus(fileErr), fileErr
	}

	req, err := deserialiseComment(responseWriter, httpRequest)
	if err != nil {
		return http.StatusBadRequest, err
	}

	req.Data.UserName = inputData.user.Username
	req.Data.CreatedTime = time.Now()

	err = inputData.store.Comments.Save(req.Data)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	createUserLeavesCommentNotification(inputData, req.Data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return errToStatus(err), err
})

func createUserLeavesCommentNotification(inputData *data, comment *comments.Comment) {

	newNotification := &notifications.Notification{
		ContextFilePath:  comment.FilePath,
		CausingUserName:  comment.UserName,
		NotificationType: "UserLeavesComment",
		CreatedTime:      time.Now(),
	}

	inputData.store.Notifications.Save(newNotification)

	usernotifications, _ := inputData.store.UserNotifications.FindByUserName(comment.UserName)

	usernotifications.AcknowledgedNotifications[newNotification.ID] = true
	_ = inputData.store.UserNotifications.Save(usernotifications)

}

// fkn errs
var commentDeleteHandler = withUser(func(responseWriter http.ResponseWriter, httpRequest *http.Request, inputData *data) (int, error) {
	commentId := strings.TrimSuffix(httpRequest.URL.Path, "/")
	commentId = strings.TrimPrefix(commentId, "/")

	reactionIdUint64, _ := strconv.ParseUint(commentId, 10, 32)
	reactionIdUint := uint(reactionIdUint64)

	comment, _ := inputData.store.Comments.GetById(reactionIdUint)
	if comment.UserName != inputData.user.Username {
		return errToStatus(errors.ErrPermissionDenied), errors.ErrPermissionDenied
	}
	//TODO: this is here to validate user has access to file, should make this consistent
	_, err := files.NewFileInfo(files.FileOptions{
		Fs:         inputData.user.Fs,
		Path:       comment.FilePath,
		Modify:     inputData.user.Perm.Modify,
		Expand:     true,
		ReadHeader: inputData.server.TypeDetectionByHeader,
		Checker:    inputData,
		Content:    true,
	})
	if err != nil {
		return errToStatus(err), err
	}

	err = inputData.store.Comments.Delete(reactionIdUint)
	if err != nil {
		return errToStatus(err), err
	}

	return http.StatusOK, nil
})
