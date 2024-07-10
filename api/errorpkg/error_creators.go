// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package errorpkg

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/kouprlabs/voltaserve/api/model"
)

func NewGroupNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"group_not_found",
		http.StatusNotFound,
		"Group not found.",
		MsgResourceNotFound,
		err,
	)
}

func NewFileNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"file_not_found",
		http.StatusNotFound,
		"File not found.",
		MsgResourceNotFound,
		err,
	)
}

func NewInvalidPathError(err error) *ErrorResponse {
	return NewErrorResponse(
		"invalid_path",
		http.StatusBadRequest,
		"Invalid path.",
		MsgInvalidRequest,
		err,
	)
}

func NewWorkspaceNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"workspace_not_found",
		http.StatusNotFound,
		"Workspace not found.",
		MsgResourceNotFound,
		err,
	)
}

func NewOrganizationNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"organization_not_found",
		http.StatusNotFound,
		"Organization not found.",
		MsgResourceNotFound,
		err,
	)
}

func NewTaskNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"task_not_found",
		http.StatusNotFound,
		"Task not found.",
		MsgResourceNotFound,
		err,
	)
}

func NewSnapshotNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"snapshot_not_found",
		http.StatusNotFound,
		"Snapshot not found.",
		"The file has no snapshots.",
		err,
	)
}

func NewS3ObjectNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"s3_object_not_found",
		http.StatusNotFound,
		"S3 object not found.",
		"The snapshot does not contain the S3 object requested.",
		err,
	)
}

func NewUserNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"user_not_found",
		http.StatusNotFound,
		"User not found.",
		MsgResourceNotFound,
		err,
	)
}

func NewInsightsNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"insights_not_found",
		http.StatusNotFound,
		"Insights not found.",
		"Insights not found.",
		err,
	)
}

func NewMosaicNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"mosaic_not_found",
		http.StatusNotFound,
		"Mosaic not found.",
		"Mosaic not found.",
		err,
	)
}

func NewWatermarkNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"watermark_not_found",
		http.StatusNotFound,
		"Watermark not found.",
		"Watermark not found.",
		err,
	)
}

func NewInvitationNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"invitation_not_found",
		http.StatusNotFound,
		"Invitation not found.",
		MsgResourceNotFound,
		err,
	)
}

func NewSnapshotLanguageNotSetError(err error) *ErrorResponse {
	return NewErrorResponse(
		"snapshot_language_not_set",
		http.StatusBadRequest,
		"Snapshot language is not set.",
		"Snapshot language is not set.",
		err,
	)
}

func NewSnapshotTextLengthExceedsLimitError(err error) *ErrorResponse {
	return NewErrorResponse(
		"snapshot_text_length_exceeds_limit",
		http.StatusBadRequest,
		"Snapshot text length exceeds limit.",
		"Snapshot text length exceeds limit.",
		err,
	)
}

func NewSnapshotCannotBePatchedError(err error) *ErrorResponse {
	return NewErrorResponse(
		"snapshot_cannot_be_patched",
		http.StatusBadRequest,
		"Snapshot cannot be patched.",
		"Snapshot cannot be patched.",
		err,
	)
}

func NewSnapshotHasPendingTaskError(err error) *ErrorResponse {
	return NewErrorResponse(
		"snapshot_has_pending_task",
		http.StatusBadRequest,
		"Snapshot has a pending task.",
		"Snapshot has a pending task.",
		err,
	)
}

func NewUnsupportedFileTypeError(err error) *ErrorResponse {
	return NewErrorResponse(
		"unsupported_file_type",
		http.StatusBadRequest,
		"File type is not supported.",
		"File type is not supported.",
		err,
	)
}

func NewTaskIsRunningError(err error) *ErrorResponse {
	return NewErrorResponse(
		"task_is_running",
		http.StatusBadRequest,
		"Task is running.",
		"Task is running.",
		err,
	)
}

func NewTaskBelongsToAnotherUserError(err error) *ErrorResponse {
	return NewErrorResponse(
		"task_belongs_to_another_user",
		http.StatusBadRequest,
		"Task belongs to another user.",
		"Task belongs to another user.",
		err,
	)
}

func NewCannotDeletePendingTaskError(err error) *ErrorResponse {
	return NewErrorResponse(
		"cannot_delete_pending_task",
		http.StatusBadRequest,
		"Cannot delete a pending task.",
		"Cannot delete a pending task.",
		err,
	)
}

func NewInternalServerError(err error) *ErrorResponse {
	return NewErrorResponse(
		"internal_server_error",
		http.StatusInternalServerError,
		"Internal server error.",
		MsgSomethingWentWrong,
		err,
	)
}

func NewOrganizationPermissionError(userID string, org model.Organization, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_organization_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"User '%s' is missing the permission '%s' for organization '%s' (%s).",
			userID, permission, org.GetName(), org.GetID(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for organization '%s'.", org.GetName()),
		nil,
	)
}

func NewCannotRemoveLastRemainingOwnerOfOrganizationError(id string) *ErrorResponse {
	return NewErrorResponse(
		"cannot_remove_last_owner_of_organization",
		http.StatusBadRequest,
		fmt.Sprintf("Cannot remove the last remaining owner of organization '%s'.", id), MsgInvalidRequest,
		nil,
	)
}

func NewGroupPermissionError(userID string, org model.Organization, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_group_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"User '%s' is missing the permission '%s' for the group '%s' (%s).",
			userID, permission, org.GetName(), org.GetID(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for the group '%s'.", org.GetName()),
		nil,
	)
}

func NewWorkspacePermissionError(userID string, workspace model.Workspace, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_workspace_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"User '%s' is missing the permission '%s' for the workspace '%s' (%s).",
			userID, permission, workspace.GetName(), workspace.GetID(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for the workspace '%s'.", workspace.GetName()),
		nil,
	)
}

func NewFilePermissionError(userID string, file model.File, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_file_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"User '%s' is missing the permission '%s' for the file '%s' (%s).",
			userID, permission, file.GetName(), file.GetID(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for the item '%s'.", file.GetName()),
		nil,
	)
}

func NewS3Error(message string) *ErrorResponse {
	return NewErrorResponse(
		"s3_error",
		http.StatusInternalServerError,
		message,
		MsgSomethingWentWrong,
		nil,
	)
}

func NewMissingQueryParamError(param string) *ErrorResponse {
	return NewErrorResponse(
		"missing_query_param",
		http.StatusBadRequest,
		fmt.Sprintf("Query param '%s' is required.", param),
		MsgInvalidRequest,
		nil,
	)
}

func NewInvalidQueryParamError(param string) *ErrorResponse {
	return NewErrorResponse(
		"invalid_query_param",
		http.StatusBadRequest,
		fmt.Sprintf("Query param '%s' is invalid.", param),
		MsgInvalidRequest,
		nil,
	)
}

func NewStorageLimitExceededError() *ErrorResponse {
	return NewErrorResponse(
		"storage_limit_exceeded",
		http.StatusForbidden,
		"Storage limit exceeded.",
		"The storage limit of your workspace has been reached, please increase it and try again.",
		nil,
	)
}

func NewInsufficientStorageCapacityError() *ErrorResponse {
	return NewErrorResponse(
		"insufficient_storage_capacity",
		http.StatusForbidden,
		"Insufficient storage capacity.",
		"The requested storage capacity is insufficient.",
		nil,
	)
}

func NewRequestBodyValidationError(err error) *ErrorResponse {
	var fields []string
	for _, e := range err.(validator.ValidationErrors) {
		fields = append(fields, e.Field())
	}
	return NewErrorResponse(
		"request_validation_error",
		http.StatusBadRequest,
		fmt.Sprintf("Failed validation for the following fields: %s.", strings.Join(fields, ",")),
		MsgInvalidRequest,
		err,
	)
}

func NewFileAlreadyChildOfDestinationError(source model.File, target model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_already_child_of_destination",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) is already a child of '%s' (%s).", source.GetName(), source.GetID(), target.GetName(), target.GetID()),
		fmt.Sprintf("Item '%s' is already within '%s'.", source.GetName(), target.GetName()),
		nil,
	)
}

func NewFileCannotBeMovedIntoItselfError(source model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_moved_into_itself",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) cannot be moved into itself.", source.GetName(), source.GetID()),
		fmt.Sprintf("Item '%s' cannot be moved into itself.", source.GetName()),
		nil,
	)
}

func NewFileIsNotAFolderError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_is_not_a_folder",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) is not a folder.", file.GetName(), file.GetID()),
		fmt.Sprintf("Item '%s' is not a folder.", file.GetName()),
		nil,
	)
}

func NewFileIsNotAFileError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_is_not_a_file",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) is not a file.", file.GetName(), file.GetID()),
		fmt.Sprintf("Item '%s' is not a file.", file.GetName()),
		nil,
	)
}

func NewTargetIsGrandChildOfSourceError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"target_is_grant_child_of_source",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) cannot be moved in another file within its own tree.", file.GetName(), file.GetID()),
		fmt.Sprintf("Item '%s' cannot be moved in another item within its own tree.", file.GetName()),
		nil,
	)
}

func NewCannotDeleteWorkspaceRootError(file model.File, workspace model.Workspace) *ErrorResponse {
	return NewErrorResponse(
		"cannot_delete_workspace_root",
		http.StatusForbidden,
		fmt.Sprintf("Cannot delete the root file (%s) of the workspace '%s' (%s).", file.GetID(), workspace.GetName(), workspace.GetID()),
		fmt.Sprintf("Cannot delete the root item of the workspace '%s'.", workspace.GetName()),
		nil,
	)
}

func NewFileCannotBeCopiedIntoOwnSubtreeError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_coped_into_own_subtree",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) cannot be copied in another file within its own subtree.", file.GetName(), file.GetID()),
		fmt.Sprintf("Item '%s' cannot be copied in another item within its own subtree.", file.GetName()),
		nil,
	)
}

func NewFileCannotBeCopiedIntoIselfError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_copied_into_itself",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) cannot be copied into itself.", file.GetName(), file.GetID()),
		fmt.Sprintf("Item '%s' cannot be copied into itself.", file.GetName()),
		nil,
	)
}

func NewFileWithSimilarNameExistsError() *ErrorResponse {
	return NewErrorResponse(
		"file_with_similar_name_exists",
		http.StatusForbidden,
		"File with similar name exists.",
		"Item with similar name exists.",
		nil,
	)
}

func NewInvalidPageParameterError() *ErrorResponse {
	return NewErrorResponse(
		"invalid_page_parameter",
		http.StatusBadRequest,
		"Invalid page parameter, must be >= 1.",
		MsgInvalidRequest,
		nil,
	)
}

func NewInvalidSizeParameterError() *ErrorResponse {
	return NewErrorResponse(
		"invalid_size_parameter",
		http.StatusBadRequest,
		"Invalid size parameter, must be >= 1.",
		MsgInvalidRequest,
		nil,
	)
}

func NewCannotAcceptNonPendingInvitationError(invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"cannot_accept_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot accept an invitation which is not pending, the status of the invitation (%s) is (%s).", invitation.GetID(), invitation.GetStatus()),
		"Cannot accept an invitation which is not pending.",
		nil,
	)
}

func NewCannotDeclineNonPendingInvitationError(invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"cannot_decline_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot decline an invitation which is not pending, the status of the invitation (%s) is (%s).", invitation.GetID(), invitation.GetStatus()),
		"Cannot decline an invitation which is not pending.",
		nil,
	)
}

func NewCannotResendNonPendingInvitationError(invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"cannot_resend_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot resend an invitation which is not pending, the status of the invitation (%s) is (%s).", invitation.GetID(), invitation.GetStatus()),
		"Cannot resend an invitation which is not pending.",
		nil,
	)
}

func NewUserNotAllowedToAcceptInvitationError(user model.User, invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_accept_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' (%s) is not allowed to accept the invitation (%s).", user.GetUsername(), user.GetID(), invitation.GetID()),
		"Not allowed to accept this invitation.",
		nil,
	)
}

func NewUserNotAllowedToDeclineInvitationError(user model.User, invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_decline_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' (%s) is not allowed to decline the invitation (%s).", user.GetUsername(), user.GetID(), invitation.GetID()),
		"Not allowed to decline this invitation.",
		nil,
	)
}

func NewUserNotAllowedToDeleteInvitationError(user model.User, invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_delete_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' (%s) not allowed to delete the invitation (%s).", user.GetUsername(), user.GetID(), invitation.GetID()),
		"Not allowed to delete this invitation.",
		nil,
	)
}

func NewUserAlreadyMemberOfOrganizationError(user model.User, org model.Organization) *ErrorResponse {
	return NewErrorResponse(
		"user_already_member_of_organization",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' (%s) is already a member of the organization '%s' (%s).", user.GetUsername(), user.GetID(), org.GetName(), org.GetID()),
		fmt.Sprintf("You are already a member of the organization '%s'.", org.GetName()),
		nil,
	)
}

func NewInvalidAPIKeyError() *ErrorResponse {
	return NewErrorResponse(
		"invalid_api_key",
		http.StatusUnauthorized,
		"Invalid API key.",
		"The API key is either missing or invalid.",
		nil,
	)
}

func NewPathVariablesAndBodyParametersNotConsistent() *ErrorResponse {
	return NewErrorResponse(
		"path_variables_and_body_parameters_not_consistent",
		http.StatusUnauthorized,
		"Path variables and body parameters are not consistent.",
		"Path variables and body parameters are not consistent.",
		nil,
	)
}
