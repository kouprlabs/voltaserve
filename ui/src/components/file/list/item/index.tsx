import { ChangeEvent, MouseEvent, useEffect } from 'react'
import { useCallback, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Box, Link as ChakraLink, Checkbox, Text } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import cx from 'classnames'
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

export type ListItemProps = {
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
  // 'gray.100', 'gray.700'
  return (
    <Box
      className={cx(
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
        'hover:bg-gray-100',
        'hover:dark:bg-gray-700',
        'active:gray-200',
        'active:dark:gray-600',
        {
          'bg-gray-100': isChecked,
          'dark:bg-gray-700': isChecked,
        },
      )}
      style={{
        width: viewType === FileViewType.List ? '100%' : width,
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
        className={cx('flex', 'items-center', 'justify-center')}
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
        className={cx(
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
            className={cx('text-center', 'no-underline', {
              'hover:no-underline': isSelectionMode,
              'hover:underline': !isSelectionMode,
            })}
            noOfLines={3}
            cursor={isSelectionMode ? 'default' : 'pointer'}
            onClick={isSelectionMode ? undefined : handleFolderLinkClick}
          >
            {file.name}
          </ChakraLink>
        )}
        {file.type === 'file' && file.status === SnapshotStatus.Ready ? (
          <ChakraLink
            className={cx('text-center', 'no-underline', {
              'hover:no-underline': isSelectionMode,
              'hover:underline': !isSelectionMode,
            })}
            noOfLines={3}
            cursor={isSelectionMode ? 'default' : 'pointer'}
            onClick={isSelectionMode ? undefined : handleFileLinkClick}
          >
            {file.name}
          </ChakraLink>
        ) : null}
        {file.type === 'file' && file.status !== SnapshotStatus.Ready ? (
          <Text
            className={cx('text-center')}
            noOfLines={3}
            onClick={handleIconClick}
          >
            {file.name}
          </Text>
        ) : null}
      </div>
      <Text noOfLines={3} className={cx('text-gray-500', 'text-center')}>
        {date}
      </Text>
    </Box>
  )
}

export default ListItem
