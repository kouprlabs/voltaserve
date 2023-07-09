import { useCallback, useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Box, Center, Stack, useColorModeValue } from '@chakra-ui/react'
import { Spinner, variables } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import FileAPI from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import Copy from '@/components/file/copy'
import Create from '@/components/file/create'
import Delete from '@/components/file/delete'
import List from '@/components/file/list'
import Move from '@/components/file/move'
import OCR from '@/components/file/ocr'
import Path from '@/components/file/path'
import Rename from '@/components/file/rename'
import Sharing from '@/components/file/sharing'
import Toolbar from '@/components/file/toolbar'
import { percentageOf } from '@/helpers/percentage-of'
import { currentUpdated, listExtended } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { selectionUpdated } from '@/store/ui/files'

let isLoadingMore = false

const WorkspaceFilesPage = () => {
  const params = useParams()
  const fileId = params.fileId as string
  const dispatch = useAppDispatch()
  const list = useAppSelector((state) => state.entities.files.list)
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const iconScale = useAppSelector((state) => state.ui.files.iconScale)
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
      const result = await FileAPI.list(fileId, {
        page: list.page + 1,
        size: FileAPI.DEFAULT_PAGE_SIZE,
        sortBy,
        sortOrder,
      })
      dispatch(listExtended(result))
    } finally {
      setIsSpinnerVisible(false)
      isLoadingMore = false
    }
  }, [fileId, list, sortBy, sortOrder, dispatch])

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
        <Path />
        <Toolbar />
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
          <List scale={iconScale} />
          {isSpinnerVisible && (
            <Center w="100%" mb={variables.spacing2Xl} justifyContent="center">
              <Spinner />
            </Center>
          )}
        </Box>
      </Stack>
      <Sharing />
      <OCR />
      <Move />
      <Copy />
      <Create />
      <Delete />
      <Rename />
    </>
  )
}

export default WorkspaceFilesPage
