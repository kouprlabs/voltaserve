import { useCallback, useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Box,
  Button,
  Center,
  Stack,
  Text,
  useColorModeValue,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { IconChevronRight } from '@koupr/ui'
import { SectionSpinner } from '@koupr/ui'
import { FcFolder } from 'react-icons/fc'
import FileAPI, { File, FileType } from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'
import Path from '@/components/common/path'

type BrowseProps = {
  onChange?: (id: string) => void
}

const Browse = ({ onChange }: BrowseProps) => {
  const { id } = useParams()
  const { data: workspace } = WorkspaceAPI.useGetById(id)
  const [folders, setFolders] = useState<File[]>([])
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [loading, setLoading] = useState(false)
  const [isSpinnerVisible, setIsSpinnerVisible] = useState(false)
  const [fileId, setFileId] = useState<string>()
  const hoverColor = useColorModeValue('gray.100', 'gray.700')
  const activeColor = useColorModeValue('gray.200', 'gray.600')
  const borderColor = useColorModeValue('gray.300', 'gray.600')

  useEffect(() => {
    if (workspace) {
      setFileId(workspace.rootId)
    }
  }, [workspace])

  useEffect(() => {
    ;(async () => {
      if (fileId) {
        try {
          const timeoutId = setTimeout(() => setIsSpinnerVisible(true), 250)
          const result = await FileAPI.list(fileId, {
            page: 1,
            type: FileType.Folder,
          })
          clearTimeout(timeoutId)
          setTotalPages(result.totalPages)
          setFolders(result.data)
        } finally {
          setIsSpinnerVisible(false)
        }
      }
    })()
  }, [fileId])

  useEffect(() => {
    if (fileId) {
      onChange?.(fileId)
    }
  }, [fileId, onChange])

  const handleLoadMore = useCallback(async (fileId: string, page: number) => {
    try {
      setLoading(true)
      const result = await FileAPI.list(fileId, {
        page,
        type: FileType.Folder,
      })
      setTotalPages(result.totalPages)
      setFolders(result.data)
      setPage(page + 1)
    } finally {
      setLoading(false)
    }
  }, [])

  if (isSpinnerVisible) {
    return <SectionSpinner />
  }

  return (
    <Stack spacing={variables.spacingSm}>
      {workspace && fileId ? (
        <Path
          rootId={workspace.rootId}
          fileId={fileId}
          maxCharacters={10}
          onClick={(fileId) => setFileId(fileId)}
        />
      ) : null}
      <Stack
        spacing={0}
        borderTop="1px solid"
        borderTopColor={borderColor}
        pt={variables.spacing}
        h={{ base: '250px', xl: '400px' }}
        overflowY="scroll"
      >
        {folders.length > 0 ? (
          folders.map((f) => (
            <Stack
              key={f.id}
              direction="row"
              alignItems="center"
              spacing={variables.spacing}
              cursor="pointer"
              _hover={{ bg: hoverColor }}
              _active={{ bg: activeColor }}
              p={variables.spacingSm}
              borderRadius={variables.borderRadiusSm}
              onClick={() => setFileId(f.id)}
            >
              <FcFolder fontSize="36px" style={{ flexShrink: 0 }} />
              <Text noOfLines={1}>{f.name}</Text>
              <Box flexGrow={1} />
              <IconChevronRight />
            </Stack>
          ))
        ) : (
          <Center h="100%">
            <Text>There are no folders.</Text>
          </Center>
        )}
      </Stack>
      {totalPages > page && fileId ? (
        <Center w="100%" p={variables.spacing}>
          <Button
            onClick={() => handleLoadMore(fileId, page)}
            isLoading={loading}
          >
            Load More
          </Button>
        </Center>
      ) : null}
    </Stack>
  )
}

export default Browse
