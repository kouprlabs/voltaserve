import { MouseEvent, useCallback, useEffect, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Wrap,
  WrapItem,
  Text,
  Center,
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
import FileAPI from '@/client/api/file'
import {
  geEditorPermission,
  ltEditorPermission,
  ltOwnerPermission,
  ltViewerPermission,
} from '@/client/api/permission'
import { REFRESH_INTERVAL, swrConfig } from '@/client/options'
import downloadFile from '@/helpers/download-file'
import store from '@/store/configure-store'
import { folderUpdated, filesUpdated } from '@/store/entities/files'
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
import ItemDragOverlay from './item-drag-overlay'
import ItemDraggableDroppable from './item-draggable-droppable'

setInterval(async () => {
  const ids = store.getState().entities.files.list?.data.map((e) => e.id) || []
  if (ids.length > 0) {
    const files = await FileAPI.batchGet({ ids })
    store.dispatch(filesUpdated(files))
  }
}, REFRESH_INTERVAL)

type ListProps = {
  scale: number
}

const List = ({ scale }: ListProps) => {
  const dispatch = useAppDispatch()
  const params = useParams()
  const fileId = params.fileId as string
  const list = useAppSelector((state) => state.entities.files.list)
  const singleFile = useAppSelector((state) =>
    state.ui.files.selection.length === 1
      ? state.entities.files.list?.data.find(
          (f) => f.id === state.ui.files.selection[0],
        )
      : null,
  )
  const [activeId, setActiveId] = useState<string | null>(null)
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [menuPosition, setMenuPosition] = useState<{ x: number; y: number }>()
  const { data: folder } = FileAPI.useGetById(fileId, swrConfig())
  const { data: itemCount } = FileAPI.useGetItemCount(fileId, swrConfig())
  const activeFile = useMemo(
    () => list?.data.find((e) => e.id === activeId),
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

  useEffect(() => {
    if (folder) {
      dispatch(folderUpdated(folder))
    }
  }, [folder, dispatch])

  const handleDragStart = useCallback((event: DragStartEvent) => {
    setActiveId(event.active.id as string)
  }, [])

  const handleDragEnd = useCallback(() => {
    setActiveId(null)
  }, [])

  return (
    <>
      <DndContext
        sensors={sensors}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        {itemCount === 0 && (
          <Center w="100%" h="300px">
            <Text>There are no items.</Text>
          </Center>
        )}
        {itemCount && itemCount > 0 && list && list.data.length > 0 ? (
          <Wrap
            spacing={variables.spacing}
            overflow="hidden"
            pb={variables.spacingLg}
          >
            {list.data.map((f) => (
              <WrapItem key={f.id}>
                <ItemDraggableDroppable
                  file={f}
                  scale={scale}
                  onContextMenu={(event: MouseEvent) => {
                    setMenuPosition({ x: event.pageX, y: event.pageY })
                    setIsMenuOpen(true)
                  }}
                />
              </WrapItem>
            ))}
          </Wrap>
        ) : null}
        <ItemDragOverlay file={activeFile!} scale={scale} />
      </DndContext>
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

export default List
