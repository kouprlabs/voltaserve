package errorpkg

import (
	"fmt"
	"net/http"
	"strings"
	"voltaserve/model"

	"github.com/go-playground/validator/v10"
)

func NewGroupNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"group_not_found",
		http.StatusNotFound,
		"Group not found",
		MsgResourceNotFound,
		err,
	)
}

func NewFileNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"file_not_found",
		http.StatusNotFound,
		"File not found",
		MsgResourceNotFound,
		err,
	)
}

func NewWorkspaceNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"workspace_not_found",
		http.StatusNotFound,
		"Workspace not found",
		MsgResourceNotFound,
		err,
	)
}

func NewOrganizationNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"organization_not_found",
		http.StatusNotFound,
		"Organization not found",
		MsgResourceNotFound,
		err,
	)
}

func NewSnapshotNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"snapshot_not_found",
		http.StatusNotFound,
		"Snapshot not found",
		"The file has no snapshots",
		err,
	)
}

func NewS3ObjectNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"s3_object_not_found",
		http.StatusNotFound,
		"S3 object not found",
		"The snapshot does not contain the S3 object requested",
		err,
	)
}

func NewInvitationNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"invitation_not_found",
		http.StatusNotFound,
		"Invitation not found",
		MsgResourceNotFound,
		err,
	)
}

func NewUserNotFoundError(err error) *ErrorResponse {
	return NewErrorResponse(
		"user_not_found",
		http.StatusNotFound,
		"User not found",
		MsgResourceNotFound,
		err,
	)
}

func NewInternalServerError(err error) *ErrorResponse {
	return NewErrorResponse(
		"internal_server_error",
		http.StatusInternalServerError,
		"Internal server error",
		MsgSomethingWentWrong,
		err,
	)
}

func NewOrganizationPermissionError(user model.UserModel, org model.OrganizationModel, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_organization_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"user '%s' (%s) is missing the permission '%s' for organization '%s' (%s)",
			user.GetUsername(), user.GetId(), permission, org.GetName(), org.GetId(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for organization '%s'", org.GetName()),
		nil,
	)
}

func NewCannotRemoveLastRemainingOwnerOfOrganizationError(id string) *ErrorResponse {
	return NewErrorResponse(
		"cannot_remove_last_owner_of_organization",
		http.StatusBadRequest,
		fmt.Sprintf("Cannot remove the last remaining owner of organization '%s'", id), MsgInvalidRequest,
		nil,
	)
}

func NewGroupPermissionError(user model.UserModel, org model.OrganizationModel, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_group_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"user '%s' (%s) is missing the permission '%s' for the group '%s' (%s)",
			user.GetUsername(), user.GetId(), permission, org.GetName(), org.GetId(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for the group '%s'", org.GetName()),
		nil,
	)
}

func NewWorkspacePermissionError(user model.UserModel, workspace model.WorkspaceModel, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_workspace_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"user '%s' (%s) is missing the permission '%s' for the workspace '%s' (%s)",
			user.GetUsername(), user.GetId(), permission, workspace.GetName(), workspace.GetId(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for the workspace '%s'", workspace.GetName()),
		nil,
	)
}

func NewFilePermissionError(user model.UserModel, file model.FileModel, permission string) *ErrorResponse {
	return NewErrorResponse(
		"missing_file_permission",
		http.StatusForbidden,
		fmt.Sprintf(
			"user '%s' (%s) is missing the permission '%s' for the file '%s' (%s)",
			user.GetUsername(), user.GetId(), permission, file.GetName(), file.GetId(),
		),
		fmt.Sprintf("Sorry, you don't have enough permissions for the item '%s'", file.GetName()),
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
		fmt.Sprintf("Query param '%s' is required", param),
		MsgInvalidRequest,
		nil,
	)
}

func NewInvalidQueryParamError(param string) *ErrorResponse {
	return NewErrorResponse(
		"invalid_query_param",
		http.StatusBadRequest,
		fmt.Sprintf("Query param '%s' is invalid", param),
		MsgInvalidRequest,
		nil,
	)
}

func NewStorageLimitExceededError() *ErrorResponse {
	return NewErrorResponse(
		"storage_limit_exceeded",
		http.StatusForbidden,
		"Storage limit exceeded",
		"The storage limit of your workspace has been reached, please increase it and try again",
		nil,
	)
}

func NewInsufficientStorageCapacityError() *ErrorResponse {
	return NewErrorResponse(
		"insufficient_storage_capacity",
		http.StatusForbidden,
		"Insufficient storage capacity",
		"The requested storage capacity is insufficient",
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
		fmt.Sprintf("Failed validation for the following fields: %s", strings.Join(fields, ",")),
		MsgInvalidRequest,
		err,
	)
}

func NewFileAlreadyChildOfDestinationError(source model.FileModel, target model.FileModel) *ErrorResponse {
	return NewErrorResponse(
		"file_already_child_of_destination",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) is already a child of '%s' (%s)", source.GetName(), source.GetId(), target.GetName(), target.GetId()),
		fmt.Sprintf("Item '%s' is already within '%s'", source.GetName(), target.GetName()),
		nil,
	)
}

func NewFileCannotBeMovedIntoItselfError(source model.FileModel) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_moved_into_itself",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) cannot be moved into itself", source.GetName(), source.GetId()),
		fmt.Sprintf("Item '%s' cannot be moved into itself", source.GetName()),
		nil,
	)
}

func NewFileIsNotAFolderError(file model.FileModel) *ErrorResponse {
	return NewErrorResponse(
		"file_is_not_a_folder",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) is not a folder", file.GetName(), file.GetId()),
		fmt.Sprintf("Item '%s' is not a folder", file.GetName()),
		nil,
	)
}

func NewTargetIsGrandChildOfSourceError(file model.FileModel) *ErrorResponse {
	return NewErrorResponse(
		"target_is_grant_child_of_source",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) cannot be moved in another file within its own tree", file.GetName(), file.GetId()),
		fmt.Sprintf("Item '%s' cannot be moved in another item within its own tree", file.GetName()),
		nil,
	)
}

func NewCannotDeleteWorkspaceRootError(file model.FileModel, workspace model.WorkspaceModel) *ErrorResponse {
	return NewErrorResponse(
		"cannot_delete_workspace_root",
		http.StatusForbidden,
		fmt.Sprintf("Cannot delete the root file (%s) of the workspace '%s' (%s)", file.GetId(), workspace.GetName(), workspace.GetId()),
		fmt.Sprintf("Cannot delete the root item of the workspace '%s'", workspace.GetName()),
		nil,
	)
}

func NewFileCannotBeCopiedIntoOwnSubtreeError(file model.FileModel) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_coped_into_own_subtree",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) cannot be copied in another file within its own subtree", file.GetName(), file.GetId()),
		fmt.Sprintf("Item '%s' cannot be copied in another item within its own subtree", file.GetName()),
		nil,
	)
}

func NewFileCannotBeCopiedIntoIselfError(file model.FileModel) *ErrorResponse {
	return NewErrorResponse(
		"file_cannot_be_copied_into_itself",
		http.StatusForbidden,
		fmt.Sprintf("File '%s' (%s) cannot be copied into itself", file.GetName(), file.GetId()),
		fmt.Sprintf("Item '%s' cannot be copied into itself", file.GetName()),
		nil,
	)
}

func NewInvalidPageParameterError() *ErrorResponse {
	return NewErrorResponse(
		"invalid_page_parameter",
		http.StatusBadRequest,
		"Invalid page parameter, must be >= 1",
		MsgInvalidRequest,
		nil,
	)
}

func NewInvalidSizeParameterError() *ErrorResponse {
	return NewErrorResponse(
		"invalid_size_parameter",
		http.StatusBadRequest,
		"Invalid size parameter, must be >= 1",
		MsgInvalidRequest,
		nil,
	)
}

func NewCannotAcceptNonPendingInvitationError(invitation model.InvitationModel) *ErrorResponse {
	return NewErrorResponse(
		"cannot_accept_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot accept an invitation which is not pending, the status of the invitation (%s) is (%s)", invitation.GetId(), invitation.GetStatus()),
		"Cannot accept an invitation which is not pending",
		nil,
	)
}

func NewCannotDeclineNonPendingInvitationError(invitation model.InvitationModel) *ErrorResponse {
	return NewErrorResponse(
		"cannot_decline_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot decline an invitation which is not pending, the status of the invitation (%s) is (%s)", invitation.GetId(), invitation.GetStatus()),
		"Cannot decline an invitation which is not pending",
		nil,
	)
}

func NewCannotResendNonPendingInvitationError(invitation model.InvitationModel) *ErrorResponse {
	return NewErrorResponse(
		"cannot_resend_non_pending_invitation",
		http.StatusForbidden,
		fmt.Sprintf("Cannot resend an invitation which is not pending, the status of the invitation (%s) is (%s)", invitation.GetId(), invitation.GetStatus()),
		"Cannot resend an invitation which is not pending",
		nil,
	)
}

func NewUserNotAllowedToAcceptInvitationError(user model.UserModel, invitation model.InvitationModel) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_accept_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' (%s) is not allowed to accept the invitation (%s)", user.GetUsername(), user.GetId(), invitation.GetId()),
		"Not allowed to accept this invitation",
		nil,
	)
}

func NewUserNotAllowedToDeclineInvitationError(user model.UserModel, invitation model.InvitationModel) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_decline_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' (%s) is not allowed to decline the invitation (%s)", user.GetUsername(), user.GetId(), invitation.GetId()),
		"Not allowed to decline this invitation",
		nil,
	)
}

func NewUserNotAllowedToDeleteInvitationError(user model.UserModel, invitation model.InvitationModel) *ErrorResponse {
	return NewErrorResponse(
		"user_not_allowed_to_delete_invitation",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' (%s) not allowed to delete the invitation (%s)", user.GetUsername(), user.GetId(), invitation.GetId()),
		"Not allowed to delete this invitation",
		nil,
	)
}

func NewUserAlreadyMemberOfOrganizationError(user model.UserModel, org model.OrganizationModel) *ErrorResponse {
	return NewErrorResponse(
		"user_already_member_of_organization",
		http.StatusForbidden,
		fmt.Sprintf("User '%s' (%s) is already a member of the organization '%s' (%s)", user.GetUsername(), user.GetId(), org.GetName(), org.GetId()),
		fmt.Sprintf("You are already a member of the organization '%s'", org.GetName()),
		nil,
	)
}
