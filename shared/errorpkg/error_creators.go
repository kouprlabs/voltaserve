// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package errorpkg

import (
	"fmt"
	"github.com/kouprlabs/voltaserve/shared/model"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewGroupNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"group_not_found",
		http.StatusNotFound,
		"Group not found.",
		"Group not found.",
		err,
	)
}

func NewFileNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"file_not_found",
		http.StatusNotFound,
		"File not found.",
		"File not found.",
		err,
	)
}

func NewInvalidPathError(err error) *ErrorResponse {
	return NewErrorResponse(
		"invalid_path",
		http.StatusBadRequest,
		"Invalid path.",
		"An invalid request was sent to the server.",
		err,
	)
}

func NewWorkspaceNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"workspace_not_found",
		http.StatusNotFound,
		"Workspace not found.",
		"Workspace not found.",
		err,
	)
}

func NewOrganizationNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"organization_not_found",
		http.StatusNotFound,
		"Organization not found.",
		"Organization not found.",
		err,
	)
}

func NewTaskNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"task_not_found",
		http.StatusNotFound,
		"Task not found.",
		"Task not found.",
		err,
	)
}

func NewSnapshotNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"snapshot_not_found",
		http.StatusNotFound,
		"Snapshot not found.",
		"Snapshot not found.",
		err,
	)
}

func NewS3ObjectNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"s3_object_not_found",
		http.StatusNotFound,
		"S3 object not found.",
		"S3 object not found.",
		err,
	)
}

func NewUserNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"user_not_found",
		http.StatusNotFound,
		"User not found.",
		"User not found.",
		err,
	)
}

func NewUserNotMemberOfOrganizationError() *ErrorResponse {
	return NewErrorResponse(
		"user_not_member_of_organization",
		http.StatusNotFound,
		"User not is not a member of organization.",
		"User not is not a member of organization.",
		nil,
	)
}

func NewPictureNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"Picture_not_found",
		http.StatusNotFound,
		"Picture not found.",
		"Picture not found.",
		err,
	)
}

func NewEntitiesNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"entities_not_found",
		http.StatusNotFound,
		"Entities not found.",
		"Entities not found.",
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

func NewInvitationNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"invitation_not_found",
		http.StatusNotFound,
		"Invitation not found.",
		"Invitation not found.",
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

func NewInternalServerError(err error) *ErrorResponse {
	return NewErrorResponse(
		"internal_server_error",
		http.StatusInternalServerError,
		"Internal server error.",
		"Oops! something went wrong.",
		err,
	)
}

func NewOrganizationPermissionError(userID string, org model.Organization, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_organization_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"User '%s' is missing permission '%s' for organization '%s'.",
			userID, permission, org.GetID(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for organization '%s'.", org.GetName()),
		nil,
	)
}

func NewCannotRemoveSoleOwnerOfOrganizationError(org model.Organization) *ErrorResponse {
	return NewErrorResponse(
		"cannot_remove_sole_owner_of_organization",
		http.StatusBadRequest,
		fmt.Sprintf("Cannot remove sole owner of organization '%s'.", org.GetID()),
		fmt.Sprintf("Cannot remove sole owner of organization '%s'.", org.GetName()),
		nil,
	)
}

func NewCannotRemoveSoleOwnerOfGroupError(group model.Group) *ErrorResponse {
	return NewErrorResponse(
		"cannot_remove_sole_owner_of_group",
		http.StatusBadRequest,
		fmt.Sprintf("Cannot remove sole owner of group '%s'.", group.GetID()),
		fmt.Sprintf("Cannot remove sole owner of group '%s'.", group.GetName()),
		nil,
	)
}

func NewGroupPermissionError(userID string, org model.Group, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_group_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"User '%s' is missing permission '%s' for group '%s'.",
			userID, permission, org.GetID(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for group '%s'.", org.GetName()),
		nil,
	)
}

func NewWorkspacePermissionError(userID string, workspace model.Workspace, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_workspace_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"User '%s' is missing permission '%s' for workspace '%s'.",
			userID, permission, workspace.GetID(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for workspace '%s'.", workspace.GetName()),
		nil,
	)
}

func NewFilePermissionError(userID string, file model.File, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_file_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"User '%s' is missing permission '%s' for file '%s'.",
			userID, permission, file.GetID(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for item '%s'.", file.GetName()),
		nil,
	)
}

func NewS3Error(message string) *ErrorResponse {
	return NewErrorResponse(
		"s3_error",
		http.StatusInternalServerError,
		message,
		"Storage error occurred.",
		nil,
	)
}

func NewMissingQueryParamError(param string) *ErrorResponse {
	return NewErrorResponse(
		"missing_query_param",
		http.StatusBadRequest,
		fmt.Sprintf("Query param '%s' is required.", param),
		"An invalid request was sent to the server.",
		nil,
	)
}

func NewInvalidQueryParamError(param string) *ErrorResponse {
	return NewErrorResponse(
		"invalid_query_param",
		http.StatusBadRequest,
		fmt.Sprintf("Query param '%s' is invalid.", param),
		"An invalid request was sent to the server.",
		nil,
	)
}

func NewStorageLimitExceededError() *ErrorResponse {
	return NewErrorResponse(
		"storage_limit_exceeded",
		http.StatusForbidden,
		"Storage limit exceeded.",
		"Storage limit of your workspace has been reached, please increase it and try again.",
		nil,
	)
}

func NewInsufficientStorageCapacityError() *ErrorResponse {
	return NewErrorResponse(
		"insufficient_storage_capacity",
		http.StatusForbidden,
		"Insufficient storage capacity.",
		"Insufficient storage capacity.",

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
		fmt.Sprintf("Failed validation for fields: %s.", strings.Join(fields, ",")),
		"An invalid request was sent to the server.",
		err,
	)
}

func NewFileAlreadyChildOfDestinationError(source model.File, target model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_already_child_of_destination",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' is already a child of '%s'.", source.GetID(), target.GetID()),
		fmt.Sprintf("Item '%s' is already within '%s'.", source.GetName(), target.GetName()),
		nil,
	)
}

func NewFileCannotBeMovedIntoItselfError(source model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_moved_into_itself",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' cannot be moved into itself.", source.GetID()),
		fmt.Sprintf("Item '%s' cannot be moved into itself.", source.GetName()),
		nil,
	)
}

func NewFileIsNotAFolderError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_is_not_a_folder",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' is not a folder.", file.GetID()),
		fmt.Sprintf("Item '%s' is not a folder.", file.GetName()),
		nil,
	)
}

func NewFileIsNotAFileError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_is_not_a_file",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' is not a file.", file.GetID()),
		fmt.Sprintf("Item '%s' is not a file.", file.GetName()),
		nil,
	)
}

func NewFilePathMissingLeadingSlash() *ErrorResponse {
	return NewErrorResponse(
		"file_path_is_missing_leading_slash",
		http.StatusBadRequest,
		"File path is missing leading slash.",
		"File path is missing leading slash.",
		nil,
	)
}

func NewFilePathOfTypeFileHasTrailingSlash() *ErrorResponse {
	return NewErrorResponse(
		"file_path_of_type_file_has_trailing_slash",
		http.StatusBadRequest,
		"File path of type file has trailing slash.",
		"File path of type file has trailing slash.",
		nil,
	)
}

func NewFileTypeIsInvalid(fileType string) *ErrorResponse {
	return NewErrorResponse(
		"file_type_is_invalid",
		http.StatusInternalServerError,
		fmt.Sprintf("File type '%s' is invalid.", fileType),
		fmt.Sprintf("File type '%s' is invalid.", fileType),
		nil,
	)
}

func NewTargetIsGrandChildOfSourceError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"target_is_grant_child_of_source",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' cannot be moved in another file within its own tree.", file.GetID()),
		fmt.Sprintf("Item '%s' cannot be moved in another item within its own tree.", file.GetName()),
		nil,
	)
}

func NewCannotDeleteWorkspaceRootError(file model.File, workspace model.Workspace) *ErrorResponse {
	return NewErrorResponse(
		"cannot_delete_workspace_root",
		http.StatusForbidden,
		fmt.Sprintf("Cannot delete root file '%s' of workspace '%s'.", file.GetID(), workspace.GetID()),
		fmt.Sprintf("Cannot delete root item of workspace '%s'.", workspace.GetName()),
		nil,
	)
}

func NewFileCannotBeCopiedIntoOwnSubtreeError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_coped_into_own_subtree",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' cannot be copied in another file within its own subtree.", file.GetID()),
		fmt.Sprintf("Item '%s' cannot be copied in another item within its own subtree.", file.GetName()),
		nil,
	)
}

func NewFileCannotBeCopiedIntoItselfError(file model.File) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_copied_into_itself",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' cannot be copied into itself.", file.GetID()),
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

func NewCannotAcceptNonPendingInvitationError(invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"cannot_accept_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot accept an invitation which is not pending, status of invitation '%s' is '%s'.", invitation.GetID(), invitation.GetStatus()),
		"Cannot accept an invitation which is not pending.",
		nil,
	)
}

func NewCannotDeclineNonPendingInvitationError(invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"cannot_decline_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot decline an invitation which is not pending, status of invitation '%s' is '%s'.", invitation.GetID(), invitation.GetStatus()),
		"Cannot decline an invitation which is not pending.",
		nil,
	)
}

func NewCannotResendNonPendingInvitationError(invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"cannot_resend_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot resend an invitation which is not pending, status of invitation '%s' is '%s'.", invitation.GetID(), invitation.GetStatus()),
		"Cannot resend an invitation which is not pending.",
		nil,
	)
}

func NewUserNotAllowedToAcceptInvitationError(user model.User, invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_accept_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' is not allowed to accept invitation '%s'.", user.GetID(), invitation.GetID()),
		"Not allowed to accept this invitation.",
		nil,
	)
}

func NewUserNotAllowedToDeclineInvitationError(user model.User, invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_decline_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' is not allowed to decline invitation '%s'.", user.GetID(), invitation.GetID()),
		"Not allowed to decline this invitation.",
		nil,
	)
}

func NewUserNotAllowedToDeleteInvitationError(user model.User, invitation model.Invitation) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_delete_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' is not allowed to delete invitation '%s'.", user.GetID(), invitation.GetID()),
		"Not allowed to delete this invitation.",
		nil,
	)
}

func NewUserAlreadyMemberOfOrganizationError(user model.User, org model.Organization) *ErrorResponse {
	return NewErrorResponse(
		"user_already_member_of_organization",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' is already a member of organization '%s'.", user.GetID(), org.GetID()),
		fmt.Sprintf("You are already a member of organization '%s'.", org.GetName()),
		nil,
	)
}

func NewInvalidAPIKeyError() *ErrorResponse {
	return NewErrorResponse(
		"invalid_api_key",
		http.StatusUnauthorized,
		"Invalid API key.",
		"API key is either missing or invalid.",
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
