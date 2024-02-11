import { ChangeEvent, MouseEvent, useEffect } from 'react'
import { useCallback, useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  Stack,
  Center,
  Link as ChakraLink,
  useColorModeValue,
  Checkbox,
  Box,
  Text,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File, SnapshotStatus } from '@/client/api/file'
import relativeDate from '@/helpers/relative-date'
import store from '@/store/configure-store'
import { useAppDispatch } from '@/store/hook'
import {
  selectionAdded,
  selectionRemoved,
  selectionUpdated,
} from '@/store/ui/files'
import { ViewType } from '@/types/file'
import Icon from './icon'
import { performMultiSelect, performRangeSelect } from './perform-select'

type ItemProps = {
  file: File
  scale: number
  viewType: ViewType
  isPresentational?: boolean
  isLoading?: boolean
  isSelectionMode?: boolean
  onContextMenu?: (event: MouseEvent) => void
}

const WIDTH = 147
const MIN_HEIGHT = 110

const Item = ({
  file,
  scale,
  viewType,
  isPresentational,
  isLoading,
  isSelectionMode,
  onContextMenu,
}: ItemProps) => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const width = useMemo(() => `${WIDTH * scale}px`, [scale])
  const minHeight = useMemo(() => `${MIN_HEIGHT * scale}px`, [scale])
  const hoverColor = useColorModeValue('gray.100', 'gray.700')
  const activeColor = useColorModeValue('gray.200', 'gray.600')
  const [isChecked, setIsChecked] = useState(false)
  const [isSelected, setIsSelected] = useState(false)
  const date = useMemo(
    () => relativeDate(new Date(file.createTime)),
    [file.createTime],
  )

  useEffect(() => {
    const unsubscribe = store.subscribe(() => {
      if (store.getState().ui.files.selection.includes(file.id)) {
        setIsSelected(true)
        setIsChecked(true)
      } else {
        if (isSelected) {
          setIsSelected(false)
        }
        if (isChecked) {
          setIsChecked(false)
        }
      }
    })
    return () => unsubscribe()
  }, [file, isSelected, isChecked])

  const handleIconClick = useCallback(
    (event: MouseEvent) => {
      event.stopPropagation()
      if (isSelectionMode) {
        setIsChecked(!isChecked)
        if (isChecked) {
          dispatch(selectionRemoved(file.id))
        } else {
          dispatch(selectionAdded(file.id))
        }
      } else {
        if (store.getState().ui.files.isMultiSelectActive) {
          performMultiSelect(file, isSelected)
        } else if (store.getState().ui.files.isRangeSelectActive) {
          performRangeSelect(file)
        } else {
          dispatch(selectionUpdated([file.id]))
        }
      }
    },
    [file, isSelected, isChecked, isSelectionMode, dispatch],
  )

  const handleFileLinkClick = useCallback(
    (event: MouseEvent) => {
      handleIconClick(event)
      window.open(`/file/${file.id}`, '_blank')?.focus()
    },
    [file.id, handleIconClick],
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
      if (!event.target.checked) {
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
        onContextMenu?.(event)
        if (!isSelected) {
          handleIconClick(event)
        }
      }
    },
    [isSelected, handleIconClick, onContextMenu],
  )

  return (
    <Stack
      direction={viewType === ViewType.List ? 'row' : 'column'}
      alignItems="center"
      position="relative"
      spacing={variables.spacingXs}
      w={viewType === ViewType.List ? '100%' : width}
      px={viewType === ViewType.List ? variables.spacing : 0}
      py={variables.spacingSm}
      _hover={{ bg: hoverColor }}
      _active={{ bg: activeColor }}
      transition="background-color 0.4s ease"
      bg={isChecked ? hoverColor : 'transparent'}
      borderRadius={variables.borderRadiusSm}
      userSelect="none"
      cursor="default"
      onClick={handleIconClick}
      onDoubleClick={handleIconDoubleClick}
      onContextMenu={handleContextMenu}
    >
      {isSelectionMode && !isPresentational ? (
        <Checkbox
          position={viewType === ViewType.List ? 'relative' : 'absolute'}
          top={viewType === ViewType.List ? 'auto' : variables.spacingSm}
          left={viewType === ViewType.List ? 'auto' : variables.spacingSm}
          isChecked={isChecked}
          zIndex={1}
          size="lg"
          onChange={handleCheckboxChange}
        />
      ) : null}
      <Center w={width} minH={minHeight}>
        <Icon file={file} scale={scale} isLoading={isLoading} />
      </Center>
      <Box
        w={width}
        title={file.name}
        px={variables.spacingXs}
        display={viewType === ViewType.List ? 'flex' : 'block'}
        flexGrow={viewType === ViewType.List ? 1 : 0}
      >
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
            onClick={handleFileLinkClick}
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
      <Text textAlign="center" noOfLines={3} color="gray.500">
        {date}
      </Text>
    </Stack>
  )
}

export default Item
