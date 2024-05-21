import {
  ChangeEvent,
  MouseEvent,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react'
import { useParams } from 'react-router-dom'
import { Portal, Menu, MenuList, MenuItem, MenuDivider } from '@chakra-ui/react'
import {
  DndContext,
  useSensors,
  PointerSensor,
  useSensor,
  DragStartEvent,
} from '@dnd-kit/core'
import cx from 'classnames'
import { FileWithPath, useDropzone } from 'react-dropzone'
import { List as ApiFileList } from '@/client/api/file'
import {
  geEditorPermission,
  ltEditorPermission,
  ltOwnerPermission,
  ltViewerPermission,
} from '@/client/api/permission'
import AnalysisOrb from '@/components/common/analysis-orb'
import downloadFile from '@/helpers/download-file'
import mapFileList from '@/helpers/map-file-list'
import {
  IconFileCopy,
  IconDownload,
  IconEdit,
  IconArrowTopRight,
  IconGroup,
  IconDelete,
  IconHistory,
  IconUpload,
} from '@/lib'
import { UploadDecorator, uploadAdded } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidOpen as aiModalDidOpen } from '@/store/ui/analysis'
import {
  copyModalDidOpen,
  deleteModalDidOpen,
  moveModalDidOpen,
  multiSelectKeyUpdated,
  rangeSelectKeyUpdated,
  renameModalDidOpen,
  selectionAdded,
  selectionUpdated,
  sharingModalDidOpen,
  snapshotListModalDidOpen,
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
  const fileUploadInput = useRef<HTMLInputElement>(null)
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
              blob: file,
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
      dispatch(selectionUpdated([]))
      window.removeEventListener('keydown', handleKeydown)
      window.removeEventListener('keyup', handleKeyup)
    }
  }, [dispatch])

  const handleDragStart = useCallback((event: DragStartEvent) => {
    dispatch(selectionAdded(event.active.id as string))
    setActiveId(event.active.id as string)
  }, [])

  const handleDragEnd = useCallback(() => {
    setActiveId(null)
  }, [])

  const handleFileChange = useCallback(
    async (event: ChangeEvent<HTMLInputElement>) => {
      const files = mapFileList(event.target.files)
      if (files.length === 1 && singleFile) {
        dispatch(
          uploadAdded(
            new UploadDecorator({
              fileId: singleFile.id,
              blob: files[0],
            }).value,
          ),
        )
        dispatch(uploadsDrawerOpened())
        if (fileUploadInput && fileUploadInput.current) {
          fileUploadInput.current.value = ''
        }
      }
    },
    [id, fileId, singleFile, dispatch],
  )

  return (
    <>
      <div
        className={cx(
          'border-2',
          { 'border-blue-500': isDragActive },
          {
            'bg-blue-100': isDragActive,
            'dark:bg-blue-950': isDragActive,
          },
          { 'border-transparent': !isDragActive },
          'rounded-md',
          'h-full',
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
              className={cx(
                'flex',
                'items-center',
                'justify-center',
                'w-full',
                'h-[300px]',
              )}
            >
              <span>There are no items.</span>
            </div>
          )}
          {viewType === FileViewType.Grid && list.totalElements > 0 ? (
            <div
              className={cx(
                'flex',
                'flex-wrap',
                'gap-1.5',
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
                    scale={scale}
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
          {viewType === FileViewType.List && list.totalElements > 0 ? (
            <div
              className={cx(
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
                    scale={scale}
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
              icon={<IconGroup />}
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
            {singleFile?.type === 'file' ? (
              <MenuItem
                icon={<IconHistory />}
                isDisabled={
                  singleFile ? ltOwnerPermission(singleFile.permission) : false
                }
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(snapshotListModalDidOpen())
                }}
              >
                Snapshots
              </MenuItem>
            ) : null}
            <MenuItem
              icon={<IconUpload />}
              isDisabled={
                !singleFile ||
                singleFile.type !== 'file' ||
                ltViewerPermission(singleFile.permission)
              }
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                const singleId = singleFile?.id
                fileUploadInput?.current?.click()
                if (singleId) {
                  dispatch(selectionUpdated([singleId]))
                }
              }}
            >
              Upload
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
            {singleFile && singleFile.type === 'file' ? (
              <>
                <MenuItem
                  icon={<AnalysisOrb width="20px" height="20px" />}
                  isDisabled={ltEditorPermission(singleFile.permission)}
                  onClick={(event: MouseEvent) => {
                    event.stopPropagation()
                    dispatch(aiModalDidOpen())
                  }}
                >
                  Analyze
                </MenuItem>
                <MenuDivider />
              </>
            ) : null}
            <MenuItem
              icon={<IconDelete />}
              className={cx('text-red-500')}
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
              icon={<IconArrowTopRight />}
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
              icon={<IconFileCopy />}
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
      <input
        ref={fileUploadInput}
        className={cx('hidden')}
        type="file"
        multiple
        onChange={handleFileChange}
      />
    </>
  )
}

export default FileList
