import { MouseEvent, useCallback, useEffect, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Wrap,
  WrapItem,
  Text,
  Portal,
  Menu,
  MenuList,
  MenuItem,
  MenuDivider,
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
import {
  DndContext,
  useSensors,
  PointerSensor,
  useSensor,
  DragStartEvent,
} from '@dnd-kit/core'
import classNames from 'classnames'
import { FileWithPath, useDropzone } from 'react-dropzone'
import { List as ApiFileList } from '@/client/api/file'
import {
  geEditorPermission,
  ltEditorPermission,
  ltOwnerPermission,
  ltViewerPermission,
} from '@/client/api/permission'
import downloadFile from '@/helpers/download-file'
import { UploadDecorator, uploadAdded } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  copyModalDidOpen,
  deleteModalDidOpen,
  moveModalDidOpen,
  multiSelectKeyUpdated,
  rangeSelectKeyUpdated,
  renameModalDidOpen,
  sharingModalDidOpen,
} from '@/store/ui/files'
import { uploadsDrawerOpened } from '@/store/ui/uploads-drawer'
import { FileViewType } from '@/types/file'
import ListDragOverlay from './list-drag-overlay'
import ListDraggableDroppable from './list-draggable-droppable'

type FileListProps = {
  list: ApiFileList
  scale: number
}

const FileList = ({ list, scale }: FileListProps) => {
  const dispatch = useAppDispatch()
  const { id, fileId } = useParams()
  const singleFile = useAppSelector((state) =>
    state.ui.files.selection.length === 1
      ? list.data.find((e) => e.id === state.ui.files.selection[0])
      : null,
  )
  const hidden = useAppSelector((state) => state.ui.files.hidden)
  const viewType = useAppSelector((state) => state.ui.files.viewType)
  const isSelectionMode = useAppSelector(
    (state) => state.ui.files.isSelectionMode,
  )
  const [activeId, setActiveId] = useState<string | null>(null)
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [menuPosition, setMenuPosition] = useState<{ x: number; y: number }>()
  const activeFile = useMemo(
    () => list.data.find((e) => e.id === activeId),
    [list, activeId],
  )
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        delay: 200,
        tolerance: 5,
      },
    }),
  )
  const onDrop = useCallback(
    (files: FileWithPath[]) => {
      if (files.length === 0) {
        return
      }
      for (const file of files) {
        dispatch(
          uploadAdded(
            new UploadDecorator({
              workspaceId: id!,
              parentId: fileId!,
              file,
            }).value,
          ),
        )
      }
      dispatch(uploadsDrawerOpened())
    },
    [id, fileId, dispatch],
  )
  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    noClick: true,
  })

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

  const handleDragStart = useCallback((event: DragStartEvent) => {
    setActiveId(event.active.id as string)
  }, [])

  const handleDragEnd = useCallback(() => {
    setActiveId(null)
  }, [])

  return (
    <>
      <div
        className={classNames(
          'border-2',
          { 'border-green-300': isDragActive },
          { 'border-transparent': !isDragActive },
          'rounded-md',
        )}
        {...getRootProps()}
      >
        <input {...getInputProps()} />
        <DndContext
          sensors={sensors}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          {list.totalElements === 0 && (
            <div
              className={classNames(
                'flex',
                'items-center',
                'justify-center',
                'w-full',
                'h-[300px]',
              )}
            >
              <Text>There are no items.</Text>
            </div>
          )}
          {viewType === FileViewType.Grid && list.totalElements > 0 ? (
            <Wrap
              spacing={variables.spacing}
              overflow="hidden"
              pb={variables.spacingLg}
            >
              {list.data
                .filter((e) => !hidden.includes(e.id))
                .map((f) => (
                  <WrapItem key={f.id}>
                    <ListDraggableDroppable
                      file={f}
                      scale={scale}
                      viewType={viewType}
                      isSelectionMode={isSelectionMode}
                      onContextMenu={(event: MouseEvent) => {
                        setMenuPosition({ x: event.pageX, y: event.pageY })
                        setIsMenuOpen(true)
                      }}
                    />
                  </WrapItem>
                ))}
            </Wrap>
          ) : null}
          {viewType === FileViewType.List && list.totalElements > 0 ? (
            <div
              className={classNames(
                'flex',
                'flex-col',
                'gap-0.5',
                'overflow-hidden',
                'pb-2.5',
              )}
            >
              {list.data
                .filter((e) => !hidden.includes(e.id))
                .map((f) => (
                  <ListDraggableDroppable
                    key={f.id}
                    file={f}
                    scale={scale * 0.5}
                    viewType={viewType}
                    isSelectionMode={isSelectionMode}
                    onContextMenu={(event: MouseEvent) => {
                      setMenuPosition({ x: event.pageX, y: event.pageY })
                      setIsMenuOpen(true)
                    }}
                  />
                ))}
            </div>
          ) : null}
          <ListDragOverlay
            file={activeFile!}
            scale={scale}
            viewType={viewType}
          />
        </DndContext>
      </div>
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
              isDisabled={
                singleFile ? ltOwnerPermission(singleFile.permission) : false
              }
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(sharingModalDidOpen())
              }}
            >
              Sharing
            </MenuItem>
            <MenuItem
              icon={<IconDownload />}
              isDisabled={
                !singleFile ||
                singleFile.type !== 'file' ||
                ltViewerPermission(singleFile.permission)
              }
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                if (singleFile) {
                  downloadFile(singleFile)
                }
              }}
            >
              Download
            </MenuItem>
            <MenuDivider />
            <MenuItem
              icon={<IconTrash />}
              color="red"
              isDisabled={
                singleFile ? ltOwnerPermission(singleFile.permission) : false
              }
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(deleteModalDidOpen())
              }}
            >
              Delete
            </MenuItem>
            <MenuItem
              icon={<IconEdit />}
              isDisabled={
                singleFile && geEditorPermission(singleFile.permission)
                  ? false
                  : true
              }
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(renameModalDidOpen())
              }}
            >
              Rename
            </MenuItem>
            <MenuItem
              icon={<IconMove />}
              isDisabled={
                singleFile ? ltEditorPermission(singleFile.permission) : false
              }
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(moveModalDidOpen())
              }}
            >
              Move
            </MenuItem>
            <MenuItem
              icon={<IconCopy />}
              isDisabled={
                singleFile ? ltEditorPermission(singleFile.permission) : false
              }
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(copyModalDidOpen())
              }}
            >
              Copy
            </MenuItem>
          </MenuList>
        </Menu>
      </Portal>
    </>
  )
}

export default FileList
