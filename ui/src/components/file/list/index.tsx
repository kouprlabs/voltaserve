// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { MouseEvent, useCallback, useEffect, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import { SectionPlaceholder } from '@koupr/ui'
import {
  DndContext,
  useSensors,
  PointerSensor,
  useSensor,
  DragStartEvent,
} from '@dnd-kit/core'
import cx from 'classnames'
import { FileWithPath, useDropzone } from 'react-dropzone'
import { useHotkeys } from 'react-hotkeys-hook'
import { FileList as ApiFileList } from '@/client/api/file'
import UploadMenu from '@/components/common/upload-menu'
import { UploadDecorator, uploadAdded } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  contextMenuDidClose,
  contextMenuDidOpen,
  copyModalDidOpen,
  deleteModalDidOpen,
  infoModalDidOpen,
  moveModalDidOpen,
  multiSelectKeyUpdated,
  rangeSelectKeyUpdated,
  renameModalDidOpen,
  selectionAdded,
  selectionUpdated,
} from '@/store/ui/files'
import { drawerDidOpen } from '@/store/ui/uploads'
import { FileViewType } from '@/types/file'
import FileMenu, { FileMenuPosition } from '../file-menu'
import ListDragOverlay from './list-drag-overlay'
import ListDraggableDroppable from './list-draggable-droppable'

type FileListProps = {
  list: ApiFileList
  scale: number
}

const FileList = ({ list, scale }: FileListProps) => {
  const dispatch = useAppDispatch()
  const { id: workspaceId, fileId } = useParams()
  const hidden = useAppSelector((state) => state.ui.files.hidden)
  const viewType = useAppSelector((state) => state.ui.files.viewType)
  const isSelectionMode = useAppSelector(
    (state) => state.ui.files.isSelectionMode,
  )
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector(
    (state) =>
      state.ui.files.isCopyModalOpen ||
      state.ui.files.isMoveModalOpen ||
      state.ui.files.isDeleteModalOpen ||
      state.ui.files.isCreateModalOpen ||
      state.ui.files.isSharingModalOpen ||
      state.ui.files.isRenameModalOpen ||
      state.ui.files.isInfoModalOpen ||
      state.ui.mosaic.isModalOpen ||
      state.ui.insights.isModalOpen,
  )
  const [activeId, setActiveId] = useState<string | null>(null)
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [menuPosition, setMenuPosition] = useState<FileMenuPosition>()
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
              workspaceId: workspaceId!,
              parentId: fileId!,
              blob: file,
            }).value,
          ),
        )
      }
      dispatch(drawerDidOpen())
    },
    [workspaceId, fileId, dispatch],
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

  useEffect(() => {
    if (isModalOpen) {
      return
    }
    if (isMenuOpen) {
      dispatch(contextMenuDidOpen())
    } else {
      dispatch(contextMenuDidClose())
    }
  }, [isMenuOpen, isModalOpen, dispatch])

  useHotkeys(
    'mod+a',
    (keyboardEvent: KeyboardEvent) => {
      if (isModalOpen) {
        return
      }
      keyboardEvent.preventDefault()
      dispatch(selectionUpdated(list?.data.map((f) => f.id)))
    },
    [list, isModalOpen, dispatch],
  )

  useHotkeys(
    'mod+c',
    () => {
      if (isModalOpen) {
        return
      }
      if (selection.length > 0) {
        dispatch(copyModalDidOpen())
      }
    },
    [selection, isModalOpen, dispatch],
  )

  useHotkeys(
    'mod+x',
    () => {
      if (isModalOpen) {
        return
      }
      if (selection.length > 0) {
        dispatch(moveModalDidOpen())
      }
    },
    [selection, isModalOpen, dispatch],
  )

  useHotkeys(
    'mod+i',
    () => {
      if (isModalOpen) {
        return
      }
      if (selection.length > 0) {
        dispatch(infoModalDidOpen())
      }
    },
    [selection, isModalOpen, dispatch],
  )

  useHotkeys(
    'mod+e, f2',
    () => {
      if (isModalOpen) {
        return
      }
      if (selection.length === 1) {
        dispatch(renameModalDidOpen())
      }
    },
    [selection, isModalOpen, dispatch],
  )

  useHotkeys(
    'backspace, delete',
    () => {
      if (isModalOpen) {
        return
      }
      if (selection.length > 0) {
        dispatch(deleteModalDidOpen())
      }
    },
    [selection, isModalOpen, dispatch],
  )

  const handleDragStart = useCallback((event: DragStartEvent) => {
    dispatch(selectionAdded(event.active.id as string))
    setActiveId(event.active.id as string)
  }, [])

  const handleDragEnd = useCallback(() => {
    setActiveId(null)
  }, [])

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
          {list.totalElements === 0 ? (
            <SectionPlaceholder
              text="There are no items."
              content={<UploadMenu />}
            />
          ) : null}
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
      <FileMenu
        isOpen={isMenuOpen}
        position={menuPosition}
        onClose={() => setIsMenuOpen(false)}
      />
    </>
  )
}

export default FileList
