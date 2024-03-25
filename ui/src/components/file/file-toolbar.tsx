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
  useToken,
} from '@chakra-ui/react'
import {
  IconAdd,
  IconCheckCircle,
  IconCircle,
  IconCopy,
  IconDotsVertical,
  IconDownload,
  IconEdit,
  IconMove,
  IconShare,
  IconTrash,
  IconUpload,
  IconRefresh,
  IconGridFill,
  IconSortUp,
  IconSortDown,
  IconCheck,
} from '@koupr/ui'
import { useSWRConfig } from 'swr'
import classNames from 'classnames'
import { FiChevronDown } from 'react-icons/fi'
import { LuGrid, LuList, LuX } from 'react-icons/lu'
import FileAPI, { List, SortBy, SortOrder } from '@/client/api/file'
import { ltEditorPermission, ltOwnerPermission } from '@/client/api/permission'
import downloadFile from '@/helpers/download-file'
import mapFileList from '@/helpers/map-file-list'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { uploadAdded, UploadDecorator } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  sortByUpdated,
  viewTypeToggled,
  selectionModeToggled,
  sortOrderToggled,
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
  const stackClassName = classNames('flex', 'flex-row', 'gap-0.5')
  const grayColor = useToken('colors', 'gray.500')

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
        return <IconCheck color="transparent" />
      }
    },
    [sortBy],
  )

  const getSortOrderIcon = useCallback(() => {
    if (sortOrder === SortOrder.Asc) {
      return <IconSortUp />
    } else if (sortOrder === SortOrder.Desc) {
      return <IconSortDown />
    }
  }, [sortOrder])

  const getViewTypeIcon = useCallback(() => {
    if (viewType === FileViewType.Grid) {
      return <LuList fontSize="16px" />
    } else if (viewType === FileViewType.List) {
      return <LuGrid fontSize="16px" />
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
              leftIcon={<FiChevronDown fontSize="16px" />}
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
              leftIcon={<IconShare />}
              onClick={() => dispatch(sharingModalDidOpen())}
            >
              Sharing
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
              leftIcon={<IconTrash />}
              color="red"
              onClick={() => dispatch(deleteModalDidOpen())}
            >
              Delete
            </Button>
          )}
          {selectionCount > 0 ? (
            <Menu>
              <MenuButton
                as={IconButton}
                icon={<IconDotsVertical />}
                variant="solid"
                aria-label=""
                isDisabled={!list}
              />
              <Portal>
                <MenuList zIndex="dropdown">
                  <MenuItem
                    icon={<IconShare />}
                    isDisabled={selectionCount === 0 || !hasOwnerPermission}
                    onClick={() => dispatch(sharingModalDidOpen())}
                  >
                    Sharing
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
                    icon={<IconTrash />}
                    color="red"
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
                    icon={<IconMove />}
                    isDisabled={selectionCount === 0 || !hasEditorPermission}
                    onClick={() => dispatch(moveModalDidOpen())}
                  >
                    Move
                  </MenuItem>
                  <MenuItem
                    icon={<IconCopy />}
                    isDisabled={selectionCount === 0 || !hasEditorPermission}
                    onClick={() => dispatch(copyModalDidOpen())}
                  >
                    Copy
                  </MenuItem>
                  {isSelectionMode ? (
                    <>
                      <MenuDivider />
                      <MenuItem
                        icon={<IconCheckCircle />}
                        onClick={handleSelectAllClick}
                      >
                        Select All
                      </MenuItem>
                      <MenuItem
                        icon={<IconCircle />}
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
        <IconButton
          icon={isSelectionMode ? <LuX /> : <IconCheckCircle />}
          isDisabled={!list}
          variant="solid"
          aria-label=""
          onClick={handleToggleSelection}
        />
        <IconButton
          icon={<IconRefresh />}
          isLoading={isRefreshing}
          isDisabled={!list}
          variant="solid"
          aria-label=""
          onClick={handleRefresh}
        />
        <Spacer />
        <div className={classNames('flex', 'flex-row', 'gap-2.5')}>
          <Slider
            w="120px"
            value={iconScale}
            min={ICON_SCALE_SLIDER_MIN}
            max={ICON_SCALE_SLIDER_MAX}
            step={ICON_SCALE_SLIDER_STEP}
            isDisabled={!list}
            onChange={handleIconScaleChange}
          >
            <SliderTrack>
              <div className={classNames('relative')} />
              <SliderFilledTrack />
            </SliderTrack>
            <SliderThumb boxSize={8}>
              <IconGridFill style={{ color: grayColor }} />
            </SliderThumb>
          </Slider>
          <div className={stackClassName}>
            <IconButton
              icon={getSortOrderIcon()}
              fontSize="16px"
              variant="solid"
              aria-label=""
              isDisabled={!list}
              onClick={handleSortOrderToggle}
            />
            <IconButton
              icon={getViewTypeIcon()}
              fontSize="16px"
              variant="solid"
              aria-label=""
              isDisabled={!list}
              onClick={handleViewTypeToggle}
            />
            <Menu>
              <MenuButton
                as={IconButton}
                icon={<IconDotsVertical />}
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
        className="hidden"
        type="file"
        multiple
        onChange={handleFileChange}
      />
      <input
        ref={folderUploadInput}
        className="hidden"
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
