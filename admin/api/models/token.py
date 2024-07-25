from .generic import GenericTokenPayload, GenericTokenResponse


class TokenPayload(GenericTokenPayload):
    pass


# class TokenRequest(GenericTokenRequest):
#     Authorization: Header


class TokenResponse(GenericTokenResponse):
    authorized: bool = False
