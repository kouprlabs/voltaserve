// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { ChangeEvent, ReactElement, useCallback, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Portal,
  Spacer,
  MenuDivider,
  Tooltip,
} from '@chakra-ui/react'
import {
  IconAdd,
  IconMoreVert,
  IconUpload,
  IconRefresh,
  IconGridView,
  IconArrowDownward,
  IconArrowUpward,
  IconCheck,
  IconLibraryAddCheck,
  IconClose,
  IconList,
  IconCloudUpload,
} from '@koupr/ui'
import cx from 'classnames'
import FileAPI, { List, SortBy, SortOrder } from '@/client/api/file'
import { ltEditorPermission } from '@/client/api/permission'
import mapFileList from '@/lib/helpers/map-file-list'
import { uploadAdded, UploadDecorator } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  sortByUpdated,
  viewTypeToggled,
  selectionModeToggled,
  sortOrderToggled,
} from '@/store/ui/files'
import {
  createModalDidOpen,
  iconScaleUpdated,
  selectionUpdated,
} from '@/store/ui/files'
import { drawerDidOpen } from '@/store/ui/uploads'
import { FileViewType } from '@/types/file'
import FileMenu from './file-menu'

export type FileToolbarProps = {
  list?: List
}

const FileToolbar = ({ list }: FileToolbarProps) => {
  const dispatch = useAppDispatch()
  const { id: workspaceId, fileId } = useParams()
  const [isRefreshing, setIsRefreshing] = useState(false)
  const count = useAppSelector(
    (state) => state.entities.files.list?.data.length,
  )
  const iconScale = useAppSelector((state) => state.ui.files.iconScale)
  const viewType = useAppSelector((state) => state.ui.files.viewType)
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const isSelectionMode = useAppSelector(
    (state) => state.ui.files.isSelectionMode,
  )
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const isContextMenuOpen = useAppSelector(
    (state) => state.ui.files.isContextMenuOpen,
  )
  const iconScales = [1, 1.25, 1.5, 1.75, 2.5]
  const fileUploadInput = useRef<HTMLInputElement>(null)
  const folderUploadInput = useRef<HTMLInputElement>(null)
  const { data: folder } = FileAPI.useGet(fileId)

  const handleUploadChange = useCallback(
    async (event: ChangeEvent<HTMLInputElement>) => {
      const files = mapFileList(event.target.files)
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
      if (fileUploadInput && fileUploadInput.current) {
        fileUploadInput.current.value = ''
      }
      if (folderUploadInput && folderUploadInput.current) {
        folderUploadInput.current.value = ''
      }
    },
    [workspaceId, fileId, dispatch],
  )

  const handleIconScaleChange = useCallback(
    (value: number) => {
      dispatch(iconScaleUpdated(value))
    },
    [dispatch],
  )

  const handleRefresh = useCallback(async () => {
    setIsRefreshing(true)
    dispatch(selectionUpdated([]))
    await mutateList?.()
    setIsRefreshing(false)
  }, [fileId, dispatch, mutateList])

  const handleSortByChange = useCallback(
    (value: SortBy) => {
      dispatch(sortByUpdated(value))
    },
    [dispatch],
  )

  const handleSortOrderToggle = useCallback(() => {
    dispatch(selectionUpdated([]))
    dispatch(sortOrderToggled())
  }, [dispatch])

  const handleViewTypeToggle = useCallback(() => {
    dispatch(selectionUpdated([]))
    dispatch(viewTypeToggled())
  }, [dispatch])

  const handleToggleSelection = useCallback(() => {
    dispatch(selectionUpdated([]))
    dispatch(selectionModeToggled())
  }, [dispatch])

  const getSortIcon = useCallback(
    (value: SortBy): ReactElement => {
      if (value === sortBy) {
        return <IconCheck />
      } else {
        return <IconCheck className={cx('text-transparent')} />
      }
    },
    [sortBy],
  )

  const getScaleIcon = useCallback(
    (value: number): ReactElement => {
      if (value === iconScale) {
        return <IconCheck />
      } else {
        return <IconCheck className={cx('text-transparent')} />
      }
    },
    [iconScale],
  )

  const getSortOrderIcon = useCallback(() => {
    if (sortOrder === SortOrder.Asc) {
      return <IconArrowUpward />
    } else if (sortOrder === SortOrder.Desc) {
      return <IconArrowDownward />
    }
  }, [sortOrder])

  const getViewTypeIcon = useCallback(() => {
    if (viewType === FileViewType.Grid) {
      return <IconList />
    } else if (viewType === FileViewType.List) {
      return <IconGridView />
    }
  }, [viewType])

  return (
    <>
      <div className={cx('flex', 'flex-row', 'gap-0.5')}>
        <Menu>
          <Tooltip label="Upload">
            <MenuButton
              as={IconButton}
              variant="solid"
              colorScheme="blue"
              icon={<IconCloudUpload />}
              isDisabled={
                !folder || ltEditorPermission(folder.permission) || !list
              }
            />
          </Tooltip>
          <Portal>
            <MenuList>
              <MenuItem
                icon={<IconAdd />}
                onClick={() => fileUploadInput?.current?.click()}
              >
                Upload Files
              </MenuItem>
              <MenuItem
                icon={<IconUpload />}
                onClick={() => folderUploadInput?.current?.click()}
              >
                Upload Folder
              </MenuItem>
            </MenuList>
          </Portal>
        </Menu>
        <Tooltip label="New Folder">
          <IconButton
            variant="outline"
            colorScheme="blue"
            icon={<IconAdd />}
            isDisabled={
              !folder || ltEditorPermission(folder.permission) || !list
            }
            onClick={() => dispatch(createModalDidOpen())}
            title="New folder"
            aria-label="New folder"
          />
        </Tooltip>
        {!isContextMenuOpen ? (
          <div className={cx('flex', 'flex-row', 'gap-0.5')}>
            <FileMenu isToolbarMode={true} />
          </div>
        ) : null}
        {count ? (
          <IconButton
            icon={isSelectionMode ? <IconClose /> : <IconLibraryAddCheck />}
            isDisabled={!list}
            variant="solid"
            title="Toggle selection mode"
            aria-label="Toggle selection mode"
            onClick={handleToggleSelection}
          />
        ) : null}
        <IconButton
          icon={<IconRefresh />}
          isLoading={isRefreshing}
          isDisabled={!list}
          variant="solid"
          title="Refresh"
          aria-label="Refresh"
          onClick={handleRefresh}
        />
        <Spacer />
        <div className={cx('flex', 'flex-row', 'gap-2.5')}>
          <div className={cx('flex', 'flex-row', 'gap-0.5')}>
            <IconButton
              icon={getSortOrderIcon()}
              variant="solid"
              title="Toggle sort order"
              aria-label="Toggle sort order"
              isDisabled={!list}
              onClick={handleSortOrderToggle}
            />
            <IconButton
              icon={getViewTypeIcon()}
              variant="solid"
              title="Toggle view type"
              aria-label="Toggle view type"
              isDisabled={!list}
              onClick={handleViewTypeToggle}
            />
            <Menu>
              <MenuButton
                as={IconButton}
                icon={<IconMoreVert />}
                variant="solid"
                title="Sort by menu"
                aria-label="Sort by menu"
                isDisabled={!list}
              />
              <Portal>
                <MenuList zIndex="dropdown">
                  <MenuItem
                    icon={getSortIcon(SortBy.Name)}
                    onClick={() => handleSortByChange(SortBy.Name)}
                  >
                    Sort By Name
                  </MenuItem>
                  <MenuItem
                    icon={getSortIcon(SortBy.Kind)}
                    onClick={() => handleSortByChange(SortBy.Kind)}
                  >
                    Sort By Kind
                  </MenuItem>
                  <MenuItem
                    icon={getSortIcon(SortBy.Size)}
                    onClick={() => handleSortByChange(SortBy.Size)}
                  >
                    Sort By Size
                  </MenuItem>
                  <MenuItem
                    icon={getSortIcon(SortBy.DateCreated)}
                    onClick={() => handleSortByChange(SortBy.DateCreated)}
                  >
                    Sort By Date Created
                  </MenuItem>
                  <MenuItem
                    icon={getSortIcon(SortBy.DateModified)}
                    onClick={() => handleSortByChange(SortBy.DateModified)}
                  >
                    Sort By Date Modified
                  </MenuItem>
                  <MenuDivider />
                  {iconScales.map((scale, index) => (
                    <MenuItem
                      key={index}
                      icon={getScaleIcon(scale)}
                      onClick={() => handleIconScaleChange(scale)}
                    >
                      {scale}x
                    </MenuItem>
                  ))}
                </MenuList>
              </Portal>
            </Menu>
          </div>
        </div>
      </div>
      <input
        ref={fileUploadInput}
        className={cx('hidden')}
        type="file"
        multiple
        onChange={handleUploadChange}
      />
      <input
        ref={folderUploadInput}
        className={cx('hidden')}
        type="file"
        /* @ts-expect-error intentionaly ignored */
        directory=""
        webkitdirectory=""
        mozdirectory=""
        onChange={handleUploadChange}
      />
    </>
  )
}

export default FileToolbar
