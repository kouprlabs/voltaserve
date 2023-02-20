import { useCallback, useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Box,
  Center,
  Spinner,
  Stack,
  useColorModeValue,
} from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import FileAPI from '@/api/file'
import { swrConfig } from '@/api/options'
import WorkspaceAPI from '@/api/workspace'
import { currentUpdated, listPatched } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { selectionUpdated } from '@/store/ui/files'
import FileCopy from '@/components/file/copy'
import FileCreate from '@/components/file/create'
import FileDelete from '@/components/file/delete'
import FileList from '@/components/file/list'
import FileMove from '@/components/file/move'
import FilePath from '@/components/file/path'
import FileRename from '@/components/file/rename'
import FileSharing from '@/components/file/sharing'
import FileToolbar from '@/components/file/toolbar'
import variables from '@/theme/variables'
import { percentageOf } from '@/helpers/percentage-of'

let isLoadingMore = false

const WorkspaceFilesPage = () => {
  const params = useParams()
  const fileId = params.fileId as string
  const dispatch = useAppDispatch()
  const list = useAppSelector((state) => state.entities.files.list)
  const [isSpinnerVisible, setIsSpinnerVisible] = useState(false)
  const borderColor = useColorModeValue('gray.300', 'gray.600')
  const listContainer = useRef<HTMLDivElement>(null)
  const { data: workspace } = WorkspaceAPI.useGetById(
    params.id as string,
    swrConfig()
  )

  useEffect(() => {
    dispatch(currentUpdated(fileId))
  }, [fileId, dispatch])

  const loadMore = useCallback(async () => {
    if (isLoadingMore || !list) {
      return
    }
    isLoadingMore = true
    setIsSpinnerVisible(true)
    try {
      const result = await FileAPI.list(
        fileId,
        FileAPI.DEFAULT_PAGE_SIZE,
        list.page + 1
      )
      dispatch(listPatched(result))
    } finally {
      setIsSpinnerVisible(false)
      isLoadingMore = false
    }
  }, [fileId, list, dispatch])

  const handleScroll = useCallback(() => {
    if (listContainer.current && list) {
      if (
        listContainer.current.offsetHeight + listContainer.current.scrollTop >=
        listContainer.current.scrollHeight -
          percentageOf(listContainer.current.scrollHeight, 50)
      ) {
        if (list.totalPages > list.page) {
          loadMore()
        }
      }
    }
  }, [loadMore, listContainer, list])

  return (
    <>
      <Helmet>{workspace && <title>{workspace.name}</title>}</Helmet>
      <Stack
        spacing={variables.spacingLg}
        w="100%"
        overflow="hidden"
        flexGrow={1}
      >
        <FilePath />
        <FileToolbar />
        <Box
          ref={listContainer}
          overflowY="auto"
          overflowX="hidden"
          borderTop="1px solid"
          borderTopColor={borderColor}
          pt={variables.spacing}
          flexGrow={1}
          onScroll={handleScroll}
          onClick={() => dispatch(selectionUpdated([]))}
        >
          <FileList />
          {isSpinnerVisible && (
            <Center w="100%" mb={variables.spacing2Xl} justifyContent="center">
              <Spinner size="sm" thickness="4px" />
            </Center>
          )}
        </Box>
      </Stack>
      <FileSharing />
      <FileMove />
      <FileCopy />
      <FileCreate />
      <FileDelete />
      <FileRename />
    </>
  )
}

export default WorkspaceFilesPage
