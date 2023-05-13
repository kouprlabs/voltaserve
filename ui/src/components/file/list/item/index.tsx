import { MouseEvent, useEffect } from 'react'
import { useCallback, useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  Stack,
  Center,
  Link as ChakraLink,
  useColorModeValue,
  Checkbox,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Box,
  MenuDivider,
  Portal,
  Text,
  VStack,
} from '@chakra-ui/react'
import {
  IconCopy,
  IconDotsVerticalSm,
  IconDownload,
  IconEdit,
  IconMove,
  IconShare,
  IconTrash,
  variables,
} from '@koupr/ui'
import { File } from '@/api/file'
import {
  ltEditorPermission,
  ltOwnerPermission,
  ltViewerPermission,
} from '@/api/permission'
import store from '@/store/configure-store'
import { useAppDispatch } from '@/store/hook'
import {
  copyModalDidOpen,
  deleteModalDidOpen,
  moveModalDidOpen,
  renameModalDidOpen,
  selectionAdded,
  selectionRemoved,
  selectionUpdated,
  sharingModalDidOpen,
} from '@/store/ui/files'
import downloadFile from '@/helpers/download-file'
import relativeDate from '@/helpers/relative-date'
import Icon from './icon'
import { performMultiSelect, performRangeSelect } from './perform-select'

export enum ItemSize {
  Normal = 'normal',
  Large = 'large',
}

type FileListItemProps = {
  file: File
  size: ItemSize
}

const FileListItem = ({ file, size }: FileListItemProps) => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const width = useMemo(() => {
    if (size === ItemSize.Normal) {
      return '150px'
    }
    if (size === ItemSize.Large) {
      return '250px'
    }
  }, [size])
  const height = useMemo(() => {
    if (size === ItemSize.Normal) {
      return '110px'
    }
    if (size === ItemSize.Large) {
      return '210px'
    }
  }, [size])
  const hoverColor = useColorModeValue('gray.100', 'gray.700')
  const activeColor = useColorModeValue('gray.200', 'gray.600')
  const [isCheckboxVisible, setIsCheckboxVisible] = useState(false)
  const [isChecked, setIsChecked] = useState(false)
  const [isSelected, setIsSelected] = useState(false)
  const date = useMemo(
    () => relativeDate(new Date(file.createTime)),
    [file.createTime]
  )

  useEffect(() => {
    const unsubscribe = store.subscribe(() => {
      if (store.getState().ui.files.selection.includes(file.id)) {
        setIsSelected(true)
        setIsChecked(true)
        setIsCheckboxVisible(true)
      } else {
        setIsSelected(false)
        setIsChecked(false)
        setIsCheckboxVisible(false)
      }
    })
    return () => unsubscribe()
  }, [file.id])

  const handleDoubleDlick = useCallback(() => {
    dispatch(selectionUpdated([]))
    if (file.type === 'folder') {
      navigate(`/workspace/${file.workspaceId}/file/${file.id}`)
    } else if (file.type === 'file') {
      window.open(`/file/${file.id}`, '_blank')?.focus()
    }
  }, [file, navigate, dispatch])

  const handleSelectionClick = useCallback(
    (event?: MouseEvent) => {
      event?.stopPropagation()
      if (store.getState().ui.files.isMultiSelectActive) {
        performMultiSelect(file, isSelected)
      } else if (store.getState().ui.files.isRangeSelectActive) {
        performRangeSelect(file)
      } else {
        dispatch(selectionUpdated([file.id]))
      }
    },
    [file, isSelected, dispatch]
  )

  return (
    <Stack
      position="relative"
      spacing={variables.spacingXs}
      pb={variables.spacingXs}
      _hover={{ bg: hoverColor }}
      _active={{ bg: activeColor }}
      transition="background-color 0.4s ease"
      bg={isChecked ? hoverColor : 'transparent'}
      borderRadius={variables.borderRadiusSm}
      userSelect="none"
      onMouseEnter={() => {
        setIsCheckboxVisible(true)
      }}
      onMouseLeave={() => {
        if (!isChecked) {
          setIsCheckboxVisible(false)
        }
      }}
      onClick={(e) => e.stopPropagation()}
    >
      {isCheckboxVisible || isSelected ? (
        <Checkbox
          position="absolute"
          top="10px"
          left="8px"
          isChecked={isChecked}
          zIndex={1}
          size="lg"
          onChange={(e) => {
            e.stopPropagation()
            if (e.target.checked) {
              setIsChecked(true)
              dispatch(selectionAdded(file.id))
            } else {
              setIsChecked(false)
              dispatch(selectionRemoved(file.id))
            }
          }}
        />
      ) : (
        <span />
      )}
      <Box>
        <Menu onOpen={() => handleSelectionClick()}>
          {isCheckboxVisible && (
            <MenuButton
              as={IconButton}
              position="absolute"
              top="7px"
              right="6px"
              icon={<IconDotsVerticalSm />}
              variant="solid"
              size="xs"
              zIndex={1}
              aria-label=""
            />
          )}
          <Portal>
            <MenuList zIndex="dropdown">
              <MenuItem
                icon={<IconShare />}
                isDisabled={ltOwnerPermission(file.permission)}
                onClick={() => dispatch(sharingModalDidOpen())}
              >
                Sharing
              </MenuItem>
              <MenuItem
                icon={<IconDownload />}
                isDisabled={
                  file.type !== 'file' || ltViewerPermission(file.permission)
                }
                onClick={() => downloadFile(file)}
              >
                Download
              </MenuItem>
              <MenuDivider />
              <MenuItem
                icon={<IconTrash />}
                color="red"
                isDisabled={ltOwnerPermission(file.permission)}
                onClick={() => dispatch(deleteModalDidOpen())}
              >
                Delete
              </MenuItem>
              <MenuItem
                icon={<IconEdit />}
                isDisabled={ltEditorPermission(file.permission)}
                onClick={() => dispatch(renameModalDidOpen())}
              >
                Rename
              </MenuItem>
              <MenuItem
                icon={<IconMove />}
                isDisabled={ltEditorPermission(file.permission)}
                onClick={() => dispatch(moveModalDidOpen())}
              >
                Move
              </MenuItem>
              <MenuItem
                icon={<IconCopy />}
                isDisabled={ltEditorPermission(file.permission)}
                onClick={() => dispatch(copyModalDidOpen())}
              >
                Copy
              </MenuItem>
            </MenuList>
          </Portal>
        </Menu>
      </Box>
      <Center
        w={width}
        h={height}
        onDoubleClick={handleDoubleDlick}
        onClick={handleSelectionClick}
      >
        <Icon file={file} size={size} />
      </Center>
      <Box w={width} title={file.name} px={variables.spacingXs}>
        {file.type === 'folder' && (
          <ChakraLink
            as={Link}
            to={`/workspace/${file.workspaceId}/file/${file.id}`}
            textAlign="center"
            noOfLines={1}
            textDecoration="none"
            _hover={{ textDecoration: 'underline' }}
          >
            {file.name}
          </ChakraLink>
        )}
        {file.type === 'file' && (
          <ChakraLink
            textAlign="center"
            noOfLines={1}
            textDecoration="none"
            _hover={{ textDecoration: 'underline' }}
            onClick={() => window.open(`/file/${file.id}`, '_blank')?.focus()}
          >
            {file.name}
          </ChakraLink>
        )}
      </Box>
      <VStack spacing={0}>
        <Text color="gray.500">{date}</Text>
      </VStack>
    </Stack>
  )
}

export default FileListItem
