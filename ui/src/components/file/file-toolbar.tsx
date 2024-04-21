import { ChangeEvent, ReactElement, useCallback, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Button,
  ButtonGroup,
  IconButton,
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  Portal,
  Spacer,
  Slider,
  SliderTrack,
  SliderFilledTrack,
  SliderThumb,
} from '@chakra-ui/react'
import { useSWRConfig } from 'swr'
import cx from 'classnames'
import FileAPI, { List, SortBy, SortOrder } from '@/client/api/file'
import { ltEditorPermission, ltOwnerPermission } from '@/client/api/permission'
import downloadFile from '@/helpers/download-file'
import mapFileList from '@/helpers/map-file-list'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import {
  IconAdd,
  IconFileCopy,
  IconMoreVert,
  IconDownload,
  IconEdit,
  IconArrowTopRight,
  IconGroup,
  IconDelete,
  IconUpload,
  IconRefresh,
  IconGridView,
  IconArrowDownward,
  IconArrowUpward,
  IconCheck,
  IconSelectCheckBox,
  IconCheckBoxOutlineBlank,
  IconLibraryAddCheck,
  IconExpandMore,
  IconClose,
  IconList,
  IconHistory,
} from '@/lib'
import { uploadAdded, UploadDecorator } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  sortByUpdated,
  viewTypeToggled,
  selectionModeToggled,
  sortOrderToggled,
  snapshotsModalDidOpen,
} from '@/store/ui/files'
import {
  copyModalDidOpen,
  createModalDidOpen,
  deleteModalDidOpen,
  iconScaleUpdated,
  moveModalDidOpen,
  renameModalDidOpen,
  selectionUpdated,
  sharingModalDidOpen,
} from '@/store/ui/files'
import { uploadsDrawerOpened } from '@/store/ui/uploads-drawer'
import { FileViewType } from '@/types/file'

const ICON_SCALE_SLIDER_STEP = 0.25
const ICON_SCALE_SLIDER_MIN = 1
const ICON_SCALE_SLIDER_MAX = ICON_SCALE_SLIDER_STEP * 9

export type FileToolbarProps = {
  list?: List
}

const FileToolbar = ({ list }: FileToolbarProps) => {
  const dispatch = useAppDispatch()
  const { mutate } = useSWRConfig()
  const { id, fileId } = useParams()
  const [isRefreshing, setIsRefreshing] = useState(false)
  const fileCount = useAppSelector(
    (state) => state.entities.files.list?.data.length,
  )
  const selectionCount = useAppSelector(
    (state) => state.ui.files.selection.length,
  )
  const singleFile = useAppSelector((state) =>
    state.ui.files.selection.length === 1
      ? list?.data.find((e) => e.id === state.ui.files.selection[0])
      : null,
  )
  const iconScale = useAppSelector((state) => state.ui.files.iconScale)
  const viewType = useAppSelector((state) => state.ui.files.viewType)
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const hasOwnerPermission = useAppSelector(
    (state) =>
      list?.data.findIndex(
        (f) =>
          state.ui.files.selection.findIndex(
            (s) => f.id === s && ltOwnerPermission(f.permission),
          ) !== -1,
      ) === -1,
  )
  const hasEditorPermission = useAppSelector(
    (state) =>
      list?.data.findIndex(
        (f) =>
          state.ui.files.selection.findIndex(
            (s) => f.id === s && ltEditorPermission(f.permission),
          ) !== -1,
      ) === -1,
  )
  const isSelectionMode = useAppSelector(
    (state) => state.ui.files.isSelectionMode,
  )
  const fileUploadInput = useRef<HTMLInputElement>(null)
  const folderUploadInput = useRef<HTMLInputElement>(null)
  const fileListSearchParams = useFileListSearchParams()
  const { data: folder } = FileAPI.useGetById(fileId)
  const stackClassName = cx('flex', 'flex-row', 'gap-0.5')

  const handleFileChange = useCallback(
    async (event: ChangeEvent<HTMLInputElement>) => {
      const files = mapFileList(event.target.files)
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
      if (fileUploadInput && fileUploadInput.current) {
        fileUploadInput.current.value = ''
      }
      if (folderUploadInput && folderUploadInput.current) {
        folderUploadInput.current.value = ''
      }
    },
    [id, fileId, dispatch],
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
    await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`)
    setIsRefreshing(false)
  }, [fileId, fileListSearchParams, mutate, dispatch])

  const handleSortByChange = useCallback(
    (value: SortBy) => {
      dispatch(sortByUpdated(value))
    },
    [dispatch],
  )

  const handleSortOrderToggle = useCallback(() => {
    dispatch(sortOrderToggled())
  }, [dispatch])

  const handleViewTypeToggle = useCallback(() => {
    dispatch(viewTypeToggled())
  }, [dispatch])

  const handleSelectAllClick = useCallback(() => {
    if (list?.data) {
      dispatch(selectionUpdated(list?.data.map((f) => f.id)))
    }
  }, [list?.data, dispatch])

  const handleToggleSelection = useCallback(() => {
    dispatch(selectionUpdated([]))
    dispatch(selectionModeToggled())
  }, [dispatch])

  const getSortByIcon = useCallback(
    (value: SortBy): ReactElement => {
      if (value === sortBy) {
        return <IconCheck />
      } else {
        return <IconCheck className={cx('text-transparent')} />
      }
    },
    [sortBy],
  )

  const getSortOrderIcon = useCallback(() => {
    if (sortOrder === SortOrder.Asc) {
      return <IconArrowDownward />
    } else if (sortOrder === SortOrder.Desc) {
      return <IconArrowUpward />
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
      <div className={stackClassName}>
        <ButtonGroup isAttached>
          <Menu>
            <MenuButton
              as={Button}
              variant="solid"
              colorScheme="blue"
              leftIcon={<IconExpandMore />}
              isDisabled={
                !folder || ltEditorPermission(folder.permission) || !list
              }
            >
              Upload
            </MenuButton>
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
          <Button
            variant="outline"
            colorScheme="blue"
            leftIcon={<IconAdd />}
            isDisabled={
              !folder || ltEditorPermission(folder.permission) || !list
            }
            onClick={() => dispatch(createModalDidOpen())}
          >
            New Folder
          </Button>
        </ButtonGroup>
        <div className={stackClassName}>
          {selectionCount > 0 && hasOwnerPermission && (
            <Button
              leftIcon={<IconGroup />}
              onClick={() => dispatch(sharingModalDidOpen())}
            >
              Sharing
            </Button>
          )}
          {singleFile?.type === 'file' && hasEditorPermission && (
            <Button
              leftIcon={<IconHistory />}
              onClick={() => dispatch(snapshotsModalDidOpen())}
            >
              Snapshots
            </Button>
          )}
          {singleFile?.type === 'file' && (
            <Button
              leftIcon={<IconDownload />}
              onClick={() => downloadFile(singleFile)}
            >
              Download
            </Button>
          )}
          {selectionCount > 0 && hasOwnerPermission && (
            <Button
              leftIcon={<IconDelete />}
              className={cx('text-red-500')}
              onClick={() => dispatch(deleteModalDidOpen())}
            >
              Delete
            </Button>
          )}
          {selectionCount > 0 ? (
            <Menu>
              <MenuButton
                as={IconButton}
                icon={<IconMoreVert />}
                variant="solid"
                aria-label=""
                isDisabled={!list}
              />
              <Portal>
                <MenuList zIndex="dropdown">
                  <MenuItem
                    icon={<IconGroup />}
                    isDisabled={selectionCount === 0 || !hasOwnerPermission}
                    onClick={() => dispatch(sharingModalDidOpen())}
                  >
                    Sharing
                  </MenuItem>
                  <MenuItem
                    icon={<IconHistory />}
                    isDisabled={
                      singleFile?.type !== 'file' || !hasOwnerPermission
                    }
                    onClick={() => dispatch(snapshotsModalDidOpen())}
                  >
                    Snapshots
                  </MenuItem>
                  <MenuItem
                    icon={<IconDownload />}
                    isDisabled={singleFile?.type !== 'file'}
                    onClick={() => {
                      if (singleFile) {
                        downloadFile(singleFile)
                      }
                    }}
                  >
                    Download
                  </MenuItem>
                  <MenuDivider />
                  <MenuItem
                    icon={<IconDelete />}
                    className={cx('text-red-500')}
                    isDisabled={selectionCount === 0 || !hasOwnerPermission}
                    onClick={() => dispatch(deleteModalDidOpen())}
                  >
                    Delete
                  </MenuItem>
                  <MenuItem
                    icon={<IconEdit />}
                    isDisabled={selectionCount !== 1 || !hasEditorPermission}
                    onClick={() => dispatch(renameModalDidOpen())}
                  >
                    Rename
                  </MenuItem>
                  <MenuItem
                    icon={<IconArrowTopRight />}
                    isDisabled={selectionCount === 0 || !hasEditorPermission}
                    onClick={() => dispatch(moveModalDidOpen())}
                  >
                    Move
                  </MenuItem>
                  <MenuItem
                    icon={<IconFileCopy />}
                    isDisabled={selectionCount === 0 || !hasEditorPermission}
                    onClick={() => dispatch(copyModalDidOpen())}
                  >
                    Copy
                  </MenuItem>
                  {isSelectionMode ? (
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
          ) : null}
        </div>
        {fileCount ? (
          <IconButton
            icon={isSelectionMode ? <IconClose /> : <IconLibraryAddCheck />}
            isDisabled={!list}
            variant="solid"
            aria-label=""
            onClick={handleToggleSelection}
          />
        ) : null}
        <IconButton
          icon={<IconRefresh />}
          isLoading={isRefreshing}
          isDisabled={!list}
          variant="solid"
          aria-label=""
          onClick={handleRefresh}
        />
        <Spacer />
        <div className={cx('flex', 'flex-row', 'gap-2.5')}>
          <Slider
            className={cx('w-[120px]')}
            value={iconScale}
            min={ICON_SCALE_SLIDER_MIN}
            max={ICON_SCALE_SLIDER_MAX}
            step={ICON_SCALE_SLIDER_STEP}
            isDisabled={!list}
            onChange={handleIconScaleChange}
          >
            <SliderTrack>
              <div className={cx('relative')} />
              <SliderFilledTrack />
            </SliderTrack>
            <SliderThumb boxSize={8}>
              <IconGridView className={cx('text-gray-500')} />
            </SliderThumb>
          </Slider>
          <div className={stackClassName}>
            <IconButton
              icon={getSortOrderIcon()}
              variant="solid"
              aria-label=""
              isDisabled={!list}
              onClick={handleSortOrderToggle}
            />
            <IconButton
              icon={getViewTypeIcon()}
              variant="solid"
              aria-label=""
              isDisabled={!list}
              onClick={handleViewTypeToggle}
            />
            <Menu>
              <MenuButton
                as={IconButton}
                icon={<IconMoreVert />}
                variant="solid"
                aria-label=""
                isDisabled={!list}
              />
              <Portal>
                <MenuList zIndex="dropdown">
                  <MenuItem
                    icon={getSortByIcon(SortBy.Name)}
                    onClick={() => handleSortByChange(SortBy.Name)}
                  >
                    Sort By Name
                  </MenuItem>
                  <MenuItem
                    icon={getSortByIcon(SortBy.Kind)}
                    onClick={() => handleSortByChange(SortBy.Kind)}
                  >
                    Sort By Kind
                  </MenuItem>
                  <MenuItem
                    icon={getSortByIcon(SortBy.Size)}
                    onClick={() => handleSortByChange(SortBy.Size)}
                  >
                    Sort By Size
                  </MenuItem>
                  <MenuItem
                    icon={getSortByIcon(SortBy.DateCreated)}
                    onClick={() => handleSortByChange(SortBy.DateCreated)}
                  >
                    Sort By Date Created
                  </MenuItem>
                  <MenuItem
                    icon={getSortByIcon(SortBy.DateModified)}
                    onClick={() => handleSortByChange(SortBy.DateModified)}
                  >
                    Sort By Date Modified
                  </MenuItem>
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
        onChange={handleFileChange}
      />
      <input
        ref={folderUploadInput}
        className={cx('hidden')}
        type="file"
        /* @ts-expect-error intentionaly ignored */
        directory=""
        webkitdirectory=""
        mozdirectory=""
        onChange={handleFileChange}
      />
    </>
  )
}

export default FileToolbar
