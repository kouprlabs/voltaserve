import { useEffect, useState } from 'react'
import { useParams, useSearchParams } from 'react-router-dom'
import { Wrap, WrapItem, Text, Center } from '@chakra-ui/react'
import { Spinner, variables } from '@koupr/ui'
import FileAPI, { FileList as FileListData } from '@/api/file'
import { swrConfig } from '@/api/options'
import { listUpdated, folderUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  multiSelectKeyUpdated,
  rangeSelectKeyUpdated,
  selectionUpdated,
} from '@/store/ui/files'
import { decodeQuery } from '@/helpers/query'
import FileListItem, { ItemSize } from './item'

const FileList = () => {
  const dispatch = useAppDispatch()
  const params = useParams()
  const workspaceId = params.id as string
  const fileId = params.fileId as string
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const list = useAppSelector((state) => state.entities.files.list)
  const [isLoading, setIsLoading] = useState(false)
  const { data: folder } = FileAPI.useGetById(fileId, swrConfig())
  const { data: itemCount } = FileAPI.useGetItemCount(fileId, swrConfig())

  useEffect(() => {
    const handleKeydown = (event: KeyboardEvent) => {
      if (event.metaKey || event.ctrlKey) {
        dispatch(multiSelectKeyUpdated(true))
      }
      if (event.shiftKey) {
        dispatch(rangeSelectKeyUpdated(true))
      }
    }
    const handleKeyup = () => {
      dispatch(multiSelectKeyUpdated(false))
      dispatch(rangeSelectKeyUpdated(false))
    }
    window.addEventListener('keydown', handleKeydown)
    window.addEventListener('keyup', handleKeyup)
    return () => {
      window.removeEventListener('keydown', handleKeydown)
      window.removeEventListener('keyup', handleKeyup)
    }
  }, [dispatch])

  useEffect(() => {
    if (folder) {
      dispatch(folderUpdated(folder))
    }
  }, [folder, dispatch])

  useEffect(() => {
    ;(async () => {
      setIsLoading(true)
      dispatch(selectionUpdated([]))
      try {
        let result: FileListData
        if (query) {
          result = await FileAPI.search(
            { text: query, parentId: fileId, workspaceId },
            FileAPI.DEFAULT_PAGE_SIZE,
            1
          )
        } else {
          result = await FileAPI.list(fileId, FileAPI.DEFAULT_PAGE_SIZE, 1)
        }
        dispatch(listUpdated(result))
      } finally {
        setIsLoading(false)
      }
    })()
  }, [workspaceId, fileId, query, dispatch])

  /* TODO: check if this is really needed */
  useEffect(() => {
    if (list) {
      ;(async () => {
        const itemCount = await FileAPI.getItemCount(fileId)
        if (list.data.length === 0 && list.totalPages > 1 && itemCount > 0) {
          setIsLoading(true)
          dispatch(selectionUpdated([]))
          try {
            const result = await FileAPI.list(
              fileId,
              FileAPI.DEFAULT_PAGE_SIZE,
              1
            )
            dispatch(listUpdated(result))
          } finally {
            setIsLoading(false)
          }
        }
      })()
    }
  }, [fileId, list, dispatch])

  if (isLoading || !list) {
    return (
      <Center w="100%" h="300px" p={variables.spacing}>
        <Spinner />
      </Center>
    )
  }

  return (
    <>
      {itemCount === 0 && (
        <Center w="100%" h="300px">
          <Text>There are no items.</Text>
        </Center>
      )}
      {itemCount && itemCount > 0 && list.data.length > 0 ? (
        <Wrap
          spacing={variables.spacing}
          overflow="hidden"
          pb={variables.spacing2Xl}
        >
          {list.data.map((f) => (
            <WrapItem key={f.id}>
              <FileListItem file={f} size={ItemSize.Normal} />
            </WrapItem>
          ))}
        </Wrap>
      ) : null}
    </>
  )
}

export default FileList
