import {
  ChangeEvent,
  ReactElement,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react'
import { useParams } from 'react-router-dom'
import {
  Button,
  Stack,
  ButtonGroup,
  Box,
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
import {
  variables,
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
import { List, SortBy, SortOrder } from '@/client/api/file'
import { ltEditorPermission, ltOwnerPermission } from '@/client/api/permission'
import downloadFile from '@/helpers/download-file'
import mapFileList from '@/helpers/map-file-list'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { uploadAdded, UploadDecorator } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  sortOrderUpdated,
  sortByUpdated,
  SORT_BY_KEY,
  SORT_ORDER_KEY,
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

const ICON_SCALE_KEY = 'voltaserve_file_icon_scale'
const SPACING = variables.spacingXs
const ICON_SCALE_SLIDER_STEP = 0.25
const ICON_SCALE_SLIDER_MIN = 1
const ICON_SCALE_SLIDER_MAX = ICON_SCALE_SLIDER_STEP * 9

const Toolbar = () => {
  const dispatch = useAppDispatch()
  const { mutate } = useSWRConfig()
  const params = useParams()
  const workspaceId = params.id as string
  const fileId = params.fileId as string
  const [isRefreshing, setIsRefreshing] = useState(false)
  const selectionCount = useAppSelector(
    (state) => state.ui.files.selection.length,
  )
  const singleFile = useAppSelector((state) =>
    state.ui.files.selection.length === 1
      ? state.entities.files.list?.data.find(
          (f) => f.id === state.ui.files.selection[0],
        )
      : null,
  )
  const folder = useAppSelector((state) => state.entities.files.folder)
  const files = useAppSelector((state) => state.entities.files.list?.data)
  const iconScale = useAppSelector((state) => state.ui.files.iconScale)
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const hasOwnerPermission = useAppSelector(
    (state) =>
      state.entities.files.list?.data.findIndex(
        (f) =>
          state.ui.files.selection.findIndex(
            (s) => f.id === s && ltOwnerPermission(f.permission),
          ) !== -1,
      ) === -1,
  )
  const hasEditorPermission = useAppSelector(
    (state) =>
      state.entities.files.list?.data.findIndex(
        (f) =>
          state.ui.files.selection.findIndex(
            (s) => f.id === s && ltEditorPermission(f.permission),
          ) !== -1,
      ) === -1,
  )
  const uploadHiddenInput = useRef<HTMLInputElement>(null)
  const fileListSearchParams = useFileListSearchParams()

  useEffect(() => {
    const iconScale = localStorage.getItem(ICON_SCALE_KEY)
    if (iconScale) {
      dispatch(iconScaleUpdated(JSON.parse(iconScale)))
    }
    const sortBy = localStorage.getItem(SORT_BY_KEY)
    if (sortBy) {
      dispatch(sortByUpdated(sortBy as SortBy))
    }
    const sortOrder = localStorage.getItem(SORT_ORDER_KEY)
    if (sortOrder) {
      dispatch(sortOrderUpdated(sortOrder as SortOrder))
    }
  }, [dispatch])

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
              workspaceId: workspaceId,
              parentId: fileId,
              file,
            }).value,
          ),
        )
      }
      dispatch(uploadsDrawerOpened())
      if (uploadHiddenInput && uploadHiddenInput.current) {
        uploadHiddenInput.current.value = ''
      }
    },
    [workspaceId, fileId, dispatch],
  )

  const handleIconScaleChange = useCallback(
    (value: number) => {
      localStorage.setItem(ICON_SCALE_KEY, JSON.stringify(value))
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
      localStorage.setItem(SORT_BY_KEY, value.toString())
      dispatch(sortByUpdated(value))
    },
    [dispatch],
  )

  const handleSortOrderToggle = useCallback(() => {
    const value: SortOrder =
      sortOrder === SortOrder.Asc ? SortOrder.Desc : SortOrder.Asc
    localStorage.setItem(SORT_ORDER_KEY, value.toString())
    dispatch(sortOrderUpdated(value))
  }, [sortOrder, dispatch])

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

  return (
    <>
      <Stack direction="row" spacing={SPACING}>
        <ButtonGroup isAttached>
          <Button
            variant="solid"
            colorScheme="blue"
            leftIcon={<IconUpload />}
            isDisabled={!folder || ltEditorPermission(folder.permission)}
            onClick={() => uploadHiddenInput?.current?.click()}
          >
            Upload File
          </Button>
          <Button
            variant="outline"
            colorScheme="blue"
            leftIcon={<IconAdd />}
            isDisabled={!folder || ltEditorPermission(folder.permission)}
            onClick={() => dispatch(createModalDidOpen())}
          >
            New Folder
          </Button>
        </ButtonGroup>
        <Stack direction="row" spacing={SPACING}>
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
          <Box>
            <Menu>
              <MenuButton
                as={IconButton}
                icon={<IconDotsVertical />}
                variant="solid"
                aria-label=""
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
                  <MenuDivider />
                  <MenuItem
                    icon={<IconCheckCircle />}
                    onClick={() => {
                      if (files) {
                        dispatch(selectionUpdated(files.map((f) => f.id)))
                      }
                    }}
                  >
                    Select All
                  </MenuItem>
                  <MenuItem
                    icon={<IconCircle />}
                    onClick={() => dispatch(selectionUpdated([]))}
                  >
                    Unselect All
                  </MenuItem>
                </MenuList>
              </Portal>
            </Menu>
          </Box>
        </Stack>
        <IconButton
          icon={<IconRefresh />}
          isLoading={isRefreshing}
          variant="solid"
          aria-label=""
          onClick={handleRefresh}
        />
        <Spacer />
        <Stack direction="row" spacing={variables.spacingLg}>
          <Slider
            w="120px"
            value={iconScale}
            min={ICON_SCALE_SLIDER_MIN}
            max={ICON_SCALE_SLIDER_MAX}
            step={ICON_SCALE_SLIDER_STEP}
            onChange={handleIconScaleChange}
          >
            <SliderTrack>
              <Box position="relative" />
              <SliderFilledTrack />
            </SliderTrack>
            <SliderThumb boxSize={8}>
              <Box color="gray" as={IconGridFill} />
            </SliderThumb>
          </Slider>
          <Stack direction="row" spacing={SPACING}>
            <IconButton
              icon={getSortOrderIcon()}
              fontSize="16px"
              variant="solid"
              aria-label=""
              onClick={handleSortOrderToggle}
            />
            <Box>
              <Menu>
                <MenuButton
                  as={IconButton}
                  icon={<IconDotsVertical />}
                  variant="solid"
                  aria-label=""
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
            </Box>
          </Stack>
        </Stack>
      </Stack>
      <input
        ref={uploadHiddenInput}
        className="hidden"
        type="file"
        multiple
        onChange={handleFileChange}
      />
    </>
  )
}

export default Toolbar
