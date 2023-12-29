import { ChangeEvent, MouseEvent, useEffect } from 'react'
import { useCallback, useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  Stack,
  Center,
  Link as ChakraLink,
  useColorModeValue,
  Checkbox,
  MenuItem,
  MenuList,
  Box,
  MenuDivider,
  Text,
  VStack,
  Menu,
  Portal,
} from '@chakra-ui/react'
import {
  IconCopy,
  IconDownload,
  IconEdit,
  IconMove,
  IconShare,
  IconTrash,
  variables,
} from '@koupr/ui'
import { File, SnapshotStatus } from '@/client/api/file'
import {
  ltEditorPermission,
  ltOwnerPermission,
  ltViewerPermission,
} from '@/client/api/permission'
import downloadFile from '@/helpers/download-file'
import relativeDate from '@/helpers/relative-date'
import store from '@/store/configure-store'
import { useAppDispatch, useAppSelector } from '@/store/hook'
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
import Icon from './icon'
import { performMultiSelect, performRangeSelect } from './perform-select'

type ItemProps = {
  file: File
  scale: number
  isPresentational?: boolean
  isLoading?: boolean
}

const WIDTH = 147
const MIN_HEIGHT = 110

const Item = ({ file, scale, isPresentational, isLoading }: ItemProps) => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const selectionCount = useAppSelector(
    (state) => state.ui.files.selection.length,
  )
  const width = useMemo(() => `${WIDTH * scale}px`, [scale])
  const minHeight = useMemo(() => `${MIN_HEIGHT * scale}px`, [scale])
  const hoverColor = useColorModeValue('gray.100', 'gray.700')
  const activeColor = useColorModeValue('gray.200', 'gray.600')
  const [isCheckboxVisible, setIsCheckboxVisible] = useState(false)
  const [isHovered, setIsHovered] = useState(false)
  const [isChecked, setIsChecked] = useState(false)
  const [isSelected, setIsSelected] = useState(false)
  const date = useMemo(
    () => relativeDate(new Date(file.createTime)),
    [file.createTime],
  )
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [menuPosition, setMenuPosition] = useState<{ x: number; y: number }>()

  useEffect(() => {
    const unsubscribe = store.subscribe(() => {
      if (store.getState().ui.files.selection.includes(file.id)) {
        setIsSelected(true)
        setIsChecked(true)
        setIsCheckboxVisible(true)
      } else {
        setIsSelected(false)
        setIsChecked(false)
        if (!isHovered) {
          setIsCheckboxVisible(false)
        }
      }
    })
    return () => unsubscribe()
  }, [file.id, isHovered])

  const handleClick = useCallback(
    (event: MouseEvent) => event.stopPropagation(),
    [],
  )

  const handleMouseEnter = useCallback(() => {
    setIsCheckboxVisible(true)
    setIsHovered(true)
  }, [])

  const handleMouseLeave = useCallback(() => {
    if (!isChecked) {
      setIsCheckboxVisible(false)
    }
    setIsHovered(false)
  }, [isChecked])

  const handleIconClick = useCallback(
    (event: MouseEvent) => {
      event.stopPropagation()
      if (store.getState().ui.files.isMultiSelectActive) {
        performMultiSelect(file, isSelected)
      } else if (store.getState().ui.files.isRangeSelectActive) {
        performRangeSelect(file)
      } else {
        dispatch(selectionUpdated([file.id]))
      }
    },
    [file, isSelected, dispatch],
  )

  const handleIconDoubleClick = useCallback(() => {
    dispatch(selectionUpdated([]))
    if (file.type === 'folder') {
      navigate(`/workspace/${file.workspaceId}/file/${file.id}`)
    } else if (file.type === 'file' && file.status === SnapshotStatus.Ready) {
      window.open(`/file/${file.id}`, '_blank')?.focus()
    }
  }, [file, navigate, dispatch])

  const handleCheckboxChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      event.stopPropagation()
      if (event.target.checked) {
        setIsChecked(true)
        dispatch(selectionAdded(file.id))
      } else {
        setIsChecked(false)
        dispatch(selectionRemoved(file.id))
      }
    },
    [file.id, dispatch],
  )

  const handleContextMenu = useCallback(
    (event: MouseEvent) => {
      if (event) {
        event.preventDefault()
        setMenuPosition({ x: event.pageX, y: event.pageY })
        setIsMenuOpen(true)
        if (!isSelected) {
          handleIconClick(event)
        }
      }
    },
    [isSelected, handleIconClick],
  )

  return (
    <Stack
      position="relative"
      spacing={variables.spacingXs}
      py={variables.spacingSm}
      _hover={{ bg: hoverColor }}
      _active={{ bg: activeColor }}
      transition="background-color 0.4s ease"
      bg={isChecked ? hoverColor : 'transparent'}
      borderRadius={variables.borderRadiusSm}
      userSelect="none"
      cursor="default"
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onClick={handleClick}
      onContextMenu={handleContextMenu}
    >
      {(isCheckboxVisible || isSelected) && !isPresentational ? (
        <Checkbox
          position="absolute"
          top={variables.spacingSm}
          left={variables.spacingSm}
          isChecked={isChecked}
          zIndex={1}
          size="lg"
          onChange={handleCheckboxChange}
        />
      ) : null}
      <Center
        w={width}
        minH={minHeight}
        onDoubleClick={handleIconDoubleClick}
        onClick={handleIconClick}
      >
        <Icon file={file} scale={scale} isLoading={isLoading} />
      </Center>
      <Box w={width} title={file.name} px={variables.spacingXs}>
        {file.type === 'folder' && (
          <ChakraLink
            as={Link}
            to={`/workspace/${file.workspaceId}/file/${file.id}`}
            textAlign="center"
            noOfLines={3}
            textDecoration="none"
            _hover={{ textDecoration: 'underline' }}
          >
            {file.name}
          </ChakraLink>
        )}
        {file.type === 'file' && file.status === SnapshotStatus.Ready ? (
          <ChakraLink
            textAlign="center"
            noOfLines={3}
            textDecoration="none"
            _hover={{ textDecoration: 'underline' }}
            onClick={(event) => {
              handleIconClick(event)
              window.open(`/file/${file.id}`, '_blank')?.focus()
            }}
          >
            {file.name}
          </ChakraLink>
        ) : null}
        {file.type === 'file' && file.status !== SnapshotStatus.Ready ? (
          <Text textAlign="center" noOfLines={3} onClick={handleIconClick}>
            {file.name}
          </Text>
        ) : null}
      </Box>
      <VStack spacing={0}>
        <Text textAlign="center" noOfLines={3} color="gray.500">
          {date}
        </Text>
      </VStack>
      <Portal>
        <Menu isOpen={isMenuOpen} onClose={() => setIsMenuOpen(false)}>
          <MenuList
            zIndex="dropdown"
            style={{
              position: 'absolute',
              left: menuPosition?.x,
              top: menuPosition?.y,
            }}
          >
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
                selectionCount !== 1 ||
                file.type !== 'file' ||
                ltViewerPermission(file.permission)
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
              isDisabled={
                selectionCount !== 1 || ltEditorPermission(file.permission)
              }
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
        </Menu>
      </Portal>
    </Stack>
  )
}

export default Item
