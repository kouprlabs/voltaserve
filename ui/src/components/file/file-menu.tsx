// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { ChangeEvent, MouseEvent, useCallback, useMemo, useRef } from 'react'
import {
  IconButton,
  Kbd,
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  MenuOptionGroup,
  Portal,
} from '@chakra-ui/react'
import {
  IconArrowTopRight,
  IconCheckBoxOutlineBlank,
  IconDelete,
  IconDownload,
  IconEdit,
  IconFileCopy,
  IconGroup,
  IconHistory,
  IconInfo,
  IconModeHeat,
  IconMoreVert,
  IconSelectCheckBox,
  IconUpload,
  IconVisibility,
} from '@koupr/ui'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { geEditorPermission, geOwnerPermission, geViewerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import downloadFile from '@/lib/helpers/download-file'
import { isImage, isMicrosoftOffice, isOpenOffice, isPDF } from '@/lib/helpers/file-extension'
import mapFileList from '@/lib/helpers/map-file-list'
import { isMacOS as helperIsMacOS } from '@/lib/helpers/os'
import { UploadDecorator, uploadAdded } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  copyModalDidOpen,
  deleteModalDidOpen,
  infoModalDidOpen,
  moveModalDidOpen,
  renameModalDidOpen,
  selectionUpdated,
  sharingModalDidOpen,
} from '@/store/ui/files'
import { modalDidOpen as insightsModalDidOpen } from '@/store/ui/insights'
import { modalDidOpen as mosaicModalDidOpen } from '@/store/ui/mosaic'
import { listModalDidOpen } from '@/store/ui/snapshots'
import { drawerDidOpen } from '@/store/ui/uploads'

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

const FileMenu = ({ position, isOpen, isToolbarMode, onClose }: FileMenuProps) => {
  const dispatch = useAppDispatch()
  const list = useAppSelector((state) => state.entities.files.list)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const { data: file } = FileAPI.useGet(selection.length === 1 ? selection[0] : undefined, swrConfig())
  const isOwnerInSelection = useMemo(
    () =>
      Boolean(
        list?.data.filter((item) => selection.includes(item.id)).every((item) => geOwnerPermission(item.permission)),
      ),
    [list, selection],
  )
  const isEditorInSelection = useMemo(
    () =>
      Boolean(
        list?.data.filter((item) => selection.includes(item.id)).every((item) => geEditorPermission(item.permission)),
      ),
    [list, selection],
  )
  const isInsightsAuthorized = useMemo(
    () =>
      file?.type === 'file' &&
      !file.snapshot?.task?.isPending &&
      (isPDF(file.snapshot?.original.extension) ||
        isMicrosoftOffice(file.snapshot?.original.extension) ||
        isOpenOffice(file.snapshot?.original.extension) ||
        isImage(file.snapshot?.original.extension)) &&
      ((geViewerPermission(file.permission) && file.snapshot?.entities) || geEditorPermission(file.permission)),
    [file],
  )
  const isMosaicAuthorized = useMemo(
    () => file?.type === 'file' && !file.snapshot?.task?.isPending && isImage(file.snapshot?.original.extension),
    [file],
  )
  const isSharingAuthorized = useMemo(() => selection.length > 0 && isOwnerInSelection, [selection, isOwnerInSelection])
  const isDeleteAuthorized = useMemo(() => selection.length > 0 && isOwnerInSelection, [selection, isOwnerInSelection])
  const isMoveAuthorized = useMemo(() => selection.length > 0 && isEditorInSelection, [selection, isEditorInSelection])
  const isCopyAuthorized = useMemo(() => selection.length > 0 && isEditorInSelection, [selection, isEditorInSelection])
  const isSnapshotsAuthorized = useMemo(() => file?.type === 'file' && geOwnerPermission(file.permission), [file])
  const isUploadAuthorized = useMemo(() => file?.type === 'file' && geEditorPermission(file.permission), [file])
  const isDownloadAuthorized = useMemo(() => file?.type === 'file' && geViewerPermission(file.permission), [file])
  const isRenameAuthorized = useMemo(() => file !== undefined && geEditorPermission(file.permission), [file])
  const isInfoAuthorized = useMemo(() => file !== undefined && geViewerPermission(file.permission), [file])
  const isToolsAuthorized = useMemo(
    () => isInsightsAuthorized || isMosaicAuthorized,
    [isInsightsAuthorized, isMosaicAuthorized],
  )
  const isManagementAuthorized = useMemo(() => {
    return isSharingAuthorized || isSnapshotsAuthorized || isUploadAuthorized || isDownloadAuthorized
  }, [isSharingAuthorized, isSnapshotsAuthorized, isUploadAuthorized, isDownloadAuthorized])
  const isMacOS = useMemo(() => helperIsMacOS(), [])
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
      <Menu isOpen={isOpen} onClose={onClose}>
        {isToolbarMode ? (
          <MenuButton
            as={IconButton}
            icon={<IconMoreVert />}
            variant="solid"
            title="File menu"
            aria-label="File menu"
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
            {isToolsAuthorized ? (
              <MenuOptionGroup>
                {isInsightsAuthorized ? (
                  <MenuItem
                    icon={<IconVisibility />}
                    onClick={(event: MouseEvent) => {
                      event.stopPropagation()
                      dispatch(insightsModalDidOpen())
                    }}
                  >
                    Insights
                  </MenuItem>
                ) : null}
                {isMosaicAuthorized ? (
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
              </MenuOptionGroup>
            ) : null}
            {isToolsAuthorized ? <MenuDivider /> : null}
            {isManagementAuthorized ? (
              <MenuOptionGroup>
                {isSharingAuthorized ? (
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
                {isSnapshotsAuthorized ? (
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
                {isUploadAuthorized ? (
                  <MenuItem
                    icon={<IconUpload />}
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
                ) : null}
                {isDownloadAuthorized ? (
                  <MenuItem
                    icon={<IconDownload />}
                    onClick={async (event: MouseEvent) => {
                      event.stopPropagation()
                      if (file) {
                        await downloadFile(file)
                      }
                    }}
                  >
                    Download
                  </MenuItem>
                ) : null}
              </MenuOptionGroup>
            ) : null}
            {isManagementAuthorized ? <MenuDivider /> : null}
            <MenuOptionGroup>
              <MenuItem
                icon={<IconDelete />}
                className={cx('text-red-500')}
                isDisabled={!isDeleteAuthorized}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(deleteModalDidOpen())
                }}
              >
                <div className={cx('flex', 'flex-row', 'justify-between')}>
                  <span>Delete</span>
                  {isMacOS ? (
                    <div>
                      <Kbd>⌘</Kbd>+<Kbd>delete</Kbd>
                    </div>
                  ) : (
                    <div>
                      <Kbd>Del</Kbd>
                    </div>
                  )}
                </div>
              </MenuItem>
              <MenuItem
                icon={<IconEdit />}
                isDisabled={!isRenameAuthorized}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(renameModalDidOpen())
                }}
              >
                <div className={cx('flex', 'flex-row', 'justify-between')}>
                  <span>Rename</span>
                  {isMacOS ? (
                    <div>
                      <Kbd>⌘</Kbd>+<Kbd>E</Kbd>
                    </div>
                  ) : (
                    <div>
                      <Kbd>F2</Kbd>
                    </div>
                  )}
                </div>
              </MenuItem>
              <MenuItem
                icon={<IconArrowTopRight />}
                isDisabled={!isMoveAuthorized}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(moveModalDidOpen())
                }}
              >
                <div className={cx('flex', 'flex-row', 'justify-between')}>
                  <span>Move</span>
                  {isMacOS ? (
                    <div>
                      <Kbd>⌘</Kbd>+<Kbd>X</Kbd>
                    </div>
                  ) : (
                    <div>
                      <Kbd>Ctrl</Kbd>+<Kbd>X</Kbd>
                    </div>
                  )}
                </div>
              </MenuItem>
              <MenuItem
                icon={<IconFileCopy />}
                isDisabled={!isCopyAuthorized}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(copyModalDidOpen())
                }}
              >
                <div className={cx('flex', 'flex-row', 'justify-between')}>
                  <span>Copy</span>
                  {isMacOS ? (
                    <div>
                      <Kbd>⌘</Kbd>+<Kbd>C</Kbd>
                    </div>
                  ) : (
                    <div>
                      <Kbd>Ctrl</Kbd>+<Kbd>C</Kbd>
                    </div>
                  )}
                </div>
              </MenuItem>
            </MenuOptionGroup>
            {isToolbarMode ? (
              <MenuOptionGroup>
                <MenuDivider />
                <MenuItem icon={<IconSelectCheckBox />} onClick={handleSelectAllClick}>
                  <div className={cx('flex', 'flex-row', 'justify-between')}>
                    <span>Select All</span>
                    {isMacOS ? (
                      <div>
                        <Kbd>⌘</Kbd>+<Kbd>A</Kbd>
                      </div>
                    ) : (
                      <div>
                        <Kbd>Ctrl</Kbd>+<Kbd>A</Kbd>
                      </div>
                    )}
                  </div>
                </MenuItem>
                <MenuItem icon={<IconCheckBoxOutlineBlank />} onClick={() => dispatch(selectionUpdated([]))}>
                  Unselect All
                </MenuItem>
              </MenuOptionGroup>
            ) : null}
            <MenuOptionGroup>
              <MenuDivider />
              <MenuItem
                icon={<IconInfo />}
                isDisabled={!isInfoAuthorized}
                onClick={(event: MouseEvent) => {
                  event.stopPropagation()
                  dispatch(infoModalDidOpen())
                }}
              >
                <div className={cx('flex', 'flex-row', 'justify-between')}>
                  <span>Info</span>
                  {isMacOS ? (
                    <div>
                      <Kbd>⌘</Kbd>+<Kbd>I</Kbd>
                    </div>
                  ) : (
                    <div>
                      <Kbd>Ctrl</Kbd>+<Kbd>I</Kbd>
                    </div>
                  )}
                </div>
              </MenuItem>
            </MenuOptionGroup>
          </MenuList>
        </Portal>
      </Menu>
      <input ref={uploadInputRef} className={cx('hidden')} type="file" multiple onChange={handleUploadInputChange} />
    </>
  )
}

export default FileMenu
