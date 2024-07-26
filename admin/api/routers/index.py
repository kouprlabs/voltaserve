from typing import Annotated

from fastapi import APIRouter, Depends, status

from ..database import fetch_index, fetch_indexes
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, IndexListResponse, IndexListRequest, IndexResponse, IndexRequest

index_api_router = APIRouter(
    prefix='/index',
    tags=['index'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)


# --- GET --- #
@index_api_router.get(path="/",
                      responses={
                          status.HTTP_200_OK: {
                              'model': IndexResponse
                          }}
                      )
async def get_index(data: Annotated[IndexRequest, Depends()]):
    index = fetch_index(_id=data.id)
    if index is None:
        raise GenericNotFoundException(detail=f'Index with id={data.id} does not exist')

    return IndexResponse(**index)


@index_api_router.get(path="/all",
                      responses={
                          status.HTTP_200_OK: {
                              'model': IndexListResponse
                          }
                      }
                      )
async def get_all_indexs(data: Annotated[IndexListRequest, Depends()]):
    indexes = fetch_indexes(page=data.page, size=data.size)
    if indexes is None:
        raise GenericNotFoundException(detail='This instance has no indexes')

    return IndexListResponse(indexes=indexes)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
