import { ChangeEvent, MouseEvent, useCallback, useMemo, useRef } from 'react'
import {
  Button,
  IconButton,
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  Portal,
} from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import {
  geEditorPermission,
  geOwnerPermission,
  geViewerPermission,
} from '@/client/api/permission'
import { Status } from '@/client/api/snapshot'
import { swrConfig } from '@/client/options'
import {
  IconFileCopy,
  IconDownload,
  IconEdit,
  IconArrowTopRight,
  IconGroup,
  IconDelete,
  IconUpload,
  IconHistory,
  IconModeHeat,
  IconSecurity,
  IconSelectCheckBox,
  IconCheckBoxOutlineBlank,
  IconMoreVert,
} from '@/lib/components/icons'
import downloadFile from '@/lib/helpers/download-file'
import {
  isImage,
  isMicrosoftOffice,
  isOpenOffice,
  isPDF,
} from '@/lib/helpers/file-extension'
import mapFileList from '@/lib/helpers/map-file-list'
import { UploadDecorator, uploadAdded } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  copyModalDidOpen,
  deleteModalDidOpen,
  moveModalDidOpen,
  renameModalDidOpen,
  selectionUpdated,
  sharingModalDidOpen,
} from '@/store/ui/files'
import { modalDidOpen as insightsModalDidOpen } from '@/store/ui/insights'
import { modalDidOpen as mosaicModalDidOpen } from '@/store/ui/mosaic'
import { listModalDidOpen } from '@/store/ui/snapshots'
import { drawerDidOpen } from '@/store/ui/uploads'
import { modalDidOpen as watermarkModalDidOpen } from '@/store/ui/watermark'
import Orb from '../common/orb'

export type FileMenuProps = {
  isOpen?: boolean
  position?: FileMenuPosition
  isToolbarMode?: boolean
  onClose?: () => void
}

export type FileMenuPosition = {
  x: number
  y: number
}

const FileMenu = ({
  position,
  isOpen,
  isToolbarMode,
  onClose,
}: FileMenuProps) => {
  const dispatch = useAppDispatch()
  const list = useAppSelector((state) => state.entities.files.list)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const { data: file } = FileAPI.useGet(
    selection.length === 1 ? selection[0] : undefined,
    swrConfig(),
  )
  const isOwnerInSelection = useMemo(
    () =>
      list?.data
        .filter((item) => selection.includes(item.id))
        .every((item) => geOwnerPermission(item.permission))
        ? true
        : false,
    [list, selection],
  )
  const isEditorInSelection = useMemo(
    () =>
      list?.data
        .filter((item) => selection.includes(item.id))
        .every((item) => geEditorPermission(item.permission))
        ? true
        : false,
    [list, selection],
  )
  const isInsightsAuthorized = useMemo(
    () =>
      file?.type === 'file' &&
      !file.snapshot?.taskId &&
      ((geViewerPermission(file.permission) && file.snapshot?.entities) ||
        geEditorPermission(file.permission)),
    [file],
  )
  const isMosaicAuthorized = useMemo(
    () =>
      file?.type === 'file' &&
      !file.snapshot?.taskId &&
      isImage(file.snapshot?.original.extension),
    [file],
  )
  const isWatermarkAuthorized = useMemo(
    () =>
      file?.type === 'file' &&
      !file.snapshot?.taskId &&
      (isPDF(file.snapshot?.original.extension) ||
        isMicrosoftOffice(file.snapshot?.original.extension) ||
        isOpenOffice(file.snapshot?.original.extension) ||
        isImage(file.snapshot?.original.extension)) &&
      geEditorPermission(file.permission),
    [file],
  )
  const isSharingAuthorized = useMemo(
    () => selection.length > 0 && isOwnerInSelection,
    [selection, isOwnerInSelection],
  )
  const isDeleteAuthorized = useMemo(
    () => selection.length > 0 && isOwnerInSelection,
    [selection, isOwnerInSelection],
  )
  const isMoveAuthorized = useMemo(
    () => selection.length > 0 && isEditorInSelection,
    [selection, isEditorInSelection],
  )
  const isCopyAuthorized = useMemo(
    () => selection.length > 0 && isEditorInSelection,
    [selection, isEditorInSelection],
  )
  const isSnapshotsAuthorized = useMemo(
    () => file?.type === 'file' && geOwnerPermission(file.permission),
    [file],
  )
  const isUploadAuthorized = useMemo(
    () => file?.type === 'file' && geEditorPermission(file.permission),
    [file],
  )
  const isDownloadAuthorized = useMemo(
    () => file?.type === 'file' && geViewerPermission(file.permission),
    [file],
  )
  const isRenameAuthorized = useMemo(
    () => file !== undefined && geEditorPermission(file.permission),
    [file],
  )
  const isAnyFeatureAuthorized = useMemo(
    () =>
      isInsightsAuthorized ||
      isSharingAuthorized ||
      isSnapshotsAuthorized ||
      isMosaicAuthorized ||
      isWatermarkAuthorized,
    [],
  )
  const uploadInputRef = useRef<HTMLInputElement>(null)

  const handleUploadInputChange = useCallback(
    async (event: ChangeEvent<HTMLInputElement>) => {
      const files = mapFileList(event.target.files)
      if (files.length === 1 && file) {
        dispatch(
          uploadAdded(
            new UploadDecorator({
              fileId: file.id,
              blob: files[0],
            }).value,
          ),
        )
        dispatch(drawerDidOpen())
        if (uploadInputRef && uploadInputRef.current) {
          uploadInputRef.current.value = ''
        }
      }
    },
    [file, uploadInputRef, dispatch],
  )

  const handleSelectAllClick = useCallback(() => {
    if (list?.data) {
      dispatch(selectionUpdated(list?.data.map((f) => f.id)))
    }
  }, [list?.data, dispatch])

  return (
    <>
      {isToolbarMode ? (
        <>
          {isInsightsAuthorized ? (
            <Button
              leftIcon={<Orb width="20px" height="20px" />}
              onClick={() => dispatch(insightsModalDidOpen())}
            >
              Insights
            </Button>
          ) : null}
          {isSharingAuthorized ? (
            <Button
              leftIcon={<IconGroup />}
              onClick={() => dispatch(sharingModalDidOpen())}
            >
              Sharing
            </Button>
          ) : null}
          {isSnapshotsAuthorized ? (
            <Button
              leftIcon={<IconHistory />}
              onClick={() => dispatch(listModalDidOpen())}
            >
              Snapshots
            </Button>
          ) : null}
          {isMosaicAuthorized ? (
            <Button
              leftIcon={<IconModeHeat />}
              onClick={() => dispatch(mosaicModalDidOpen())}
            >
              Mosaic
            </Button>
          ) : null}
          {isWatermarkAuthorized ? (
            <Button
              leftIcon={<IconSecurity />}
              onClick={() => dispatch(watermarkModalDidOpen())}
            >
              Watermark
            </Button>
          ) : null}
        </>
      ) : null}
      <Menu isOpen={isOpen} onClose={onClose}>
        {isToolbarMode ? (
          <MenuButton
            as={IconButton}
            icon={<IconMoreVert />}
            variant="solid"
            aria-label=""
          />
        ) : null}
        <Portal>
          <MenuList
            zIndex="dropdown"
            style={
              position
                ? {
                    position: 'absolute',
                    left: position?.x,
                    top: position?.y,
                  }
                : undefined
            }
          >
            {!isToolbarMode && isInsightsAuthorized ? (
              <>
                <MenuItem
                  icon={<Orb width="16px" height="16px" />}
                  onClick={(event: MouseEvent) => {
                    event.stopPropagation()
                    dispatch(insightsModalDidOpen())
                  }}
                >
                  Insights
                </MenuItem>
              </>
            ) : null}
            {!isToolbarMode && isSharingAuthorized ? (
              <MenuItem
                icon={<IconGroup />}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(sharingModalDidOpen())
                }}
              >
                Sharing
              </MenuItem>
            ) : null}
            {!isToolbarMode && isSnapshotsAuthorized ? (
              <MenuItem
                icon={<IconHistory />}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(listModalDidOpen())
                }}
              >
                Snapshots
              </MenuItem>
            ) : null}
            {!isToolbarMode && isMosaicAuthorized ? (
              <MenuItem
                icon={<IconModeHeat />}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(mosaicModalDidOpen())
                }}
              >
                Mosaic
              </MenuItem>
            ) : null}
            {!isToolbarMode && isWatermarkAuthorized ? (
              <MenuItem
                icon={<IconSecurity />}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(watermarkModalDidOpen())
                }}
              >
                Watermark
              </MenuItem>
            ) : null}
            {isAnyFeatureAuthorized && !isToolbarMode ? <MenuDivider /> : null}
            <MenuItem
              icon={<IconUpload />}
              isDisabled={!isUploadAuthorized}
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                const singleId = file?.id
                uploadInputRef?.current?.click()
                if (singleId) {
                  dispatch(selectionUpdated([singleId]))
                }
              }}
            >
              Upload
            </MenuItem>
            <MenuItem
              icon={<IconDownload />}
              isDisabled={!isDownloadAuthorized}
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                if (file) {
                  downloadFile(file)
                }
              }}
            >
              Download
            </MenuItem>
            <MenuDivider />
            <MenuItem
              icon={<IconDelete />}
              className={cx('text-red-500')}
              isDisabled={!isDeleteAuthorized}
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(deleteModalDidOpen())
              }}
            >
              Delete
            </MenuItem>
            <MenuItem
              icon={<IconEdit />}
              isDisabled={!isRenameAuthorized}
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(renameModalDidOpen())
              }}
            >
              Rename
            </MenuItem>
            <MenuItem
              icon={<IconArrowTopRight />}
              isDisabled={!isMoveAuthorized}
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(moveModalDidOpen())
              }}
            >
              Move
            </MenuItem>
            <MenuItem
              icon={<IconFileCopy />}
              isDisabled={!isCopyAuthorized}
              onClick={(event: MouseEvent) => {
                event.stopPropagation()
                dispatch(copyModalDidOpen())
              }}
            >
              Copy
            </MenuItem>
            {isToolbarMode ? (
              <>
                <MenuDivider />
                <MenuItem
                  icon={<IconSelectCheckBox />}
                  onClick={handleSelectAllClick}
                >
                  Select All
                </MenuItem>
                <MenuItem
                  icon={<IconCheckBoxOutlineBlank />}
                  onClick={() => dispatch(selectionUpdated([]))}
                >
                  Unselect All
                </MenuItem>
              </>
            ) : null}
          </MenuList>
        </Portal>
      </Menu>
      <input
        ref={uploadInputRef}
        className={cx('hidden')}
        type="file"
        multiple
        onChange={handleUploadInputChange}
      />
    </>
  )
}

export default FileMenu
