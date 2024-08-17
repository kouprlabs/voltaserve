from fastapi import status

errors = {
    status.HTTP_500_INTERNAL_SERVER_ERROR: "internal_server_error",
    status.HTTP_400_BAD_REQUEST: "bad_request",
    status.HTTP_404_NOT_FOUND: "resource_not_found"
}
