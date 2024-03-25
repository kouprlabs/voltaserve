import { ChangeEvent, MouseEvent, useEffect } from 'react'
import { useCallback, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Box,
  Link as ChakraLink,
  Checkbox,
  Text,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import classNames from 'classnames'
import { SnapshotStatus } from '@/client/api/file'
import relativeDate from '@/helpers/relative-date'
import store from '@/store/configure-store'
import { useAppDispatch } from '@/store/hook'
import {
  selectionAdded,
  selectionRemoved,
  selectionUpdated,
} from '@/store/ui/files'
import { FileCommonProps, FileViewType } from '@/types/file'
import ItemIcon from './icon'
import { performMultiSelect, performRangeSelect } from './item-perform-select'

type ListItemProps = {
  onContextMenu?: (event: MouseEvent) => void
} & FileCommonProps

const WIDTH = 147
const MIN_HEIGHT = 110

const ListItem = ({
  file,
  scale,
  viewType,
  isPresentational,
  isLoading,
  isSelectionMode,
  onContextMenu,
}: ListItemProps) => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const [isChecked, setIsChecked] = useState(false)
  const [isSelected, setIsSelected] = useState(false)
  const date = relativeDate(new Date(file.createTime))
  const hoverColor = useToken(
    'colors',
    useColorModeValue('gray.100', 'gray.700'),
  )
  const activeColor = useToken(
    'colors',
    useColorModeValue('gray.200', 'gray.600'),
  )
  const width = `${WIDTH * scale}px`
  const minHeight = `${MIN_HEIGHT * scale}px`

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

  const handleFolderLinkClick = useCallback(() => {
    navigate(`/workspace/${file.workspaceId}/file/${file.id}`)
  }, [file.id, file.workspaceId, navigate])

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
    <Box
      className={classNames(
        'relative',
        'flex',
        { 'flex-col': viewType === FileViewType.Grid },
        { 'flex-row': viewType === FileViewType.List },
        'items-center',
        'gap-0.5',
        { 'px-1.5': viewType === FileViewType.List },
        { 'px-0': viewType === FileViewType.Grid },
        'py-1',
        'transition',
        'duration-400',
        'ease-in-out',
        { 'bg-transparent': !isChecked },
        'rounded-md',
        'select-none',
        'cursor-default',
      )}
      _hover={{ bg: hoverColor }}
      _active={{ bg: activeColor }}
      style={{
        width: viewType === FileViewType.List ? '100%' : width,
        background: isChecked ? hoverColor : undefined,
      }}
      onClick={handleIconClick}
      onDoubleClick={isSelectionMode ? undefined : handleIconDoubleClick}
      onContextMenu={isSelectionMode ? undefined : handleContextMenu}
    >
      {isSelectionMode && !isPresentational ? (
        <Checkbox
          position={viewType === FileViewType.List ? 'relative' : 'absolute'}
          top={viewType === FileViewType.List ? 'auto' : variables.spacingSm}
          left={viewType === FileViewType.List ? 'auto' : variables.spacingSm}
          isChecked={isChecked}
          zIndex={1}
          size="lg"
          onChange={handleCheckboxChange}
        />
      ) : null}
      <div
        className={classNames('flex', 'items-center', 'justify-center')}
        style={{ width, minHeight }}
      >
        <ItemIcon
          file={file}
          scale={scale}
          viewType={viewType}
          isLoading={isLoading}
        />
      </div>
      <div
        className={classNames(
          'px-0.5',
          { 'flex': viewType === FileViewType.List },
          { 'block': viewType === FileViewType.Grid },
          { 'grow': viewType === FileViewType.List },
          { 'grow-0': viewType === FileViewType.Grid },
        )}
        style={{ width }}
        title={file.name}
      >
        {file.type === 'folder' && (
          <ChakraLink
            textAlign="center"
            noOfLines={3}
            textDecoration="none"
            cursor={isSelectionMode ? 'default' : 'pointer'}
            _hover={{ textDecoration: isSelectionMode ? 'none' : 'underline' }}
            onClick={isSelectionMode ? undefined : handleFolderLinkClick}
          >
            {file.name}
          </ChakraLink>
        )}
        {file.type === 'file' && file.status === SnapshotStatus.Ready ? (
          <ChakraLink
            textAlign="center"
            noOfLines={3}
            textDecoration="none"
            cursor={isSelectionMode ? 'default' : 'pointer'}
            _hover={{ textDecoration: isSelectionMode ? 'none' : 'underline' }}
            onClick={isSelectionMode ? undefined : handleFileLinkClick}
          >
            {file.name}
          </ChakraLink>
        ) : null}
        {file.type === 'file' && file.status !== SnapshotStatus.Ready ? (
          <Text textAlign="center" noOfLines={3} onClick={handleIconClick}>
            {file.name}
          </Text>
        ) : null}
      </div>
      <Text textAlign="center" noOfLines={3} color="gray.500">
        {date}
      </Text>
    </Box>
  )
}

export default ListItem
