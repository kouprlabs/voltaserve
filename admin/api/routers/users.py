from fastapi import APIRouter, Depends, status

users_api_router = APIRouter(
    prefix='/users'
)
