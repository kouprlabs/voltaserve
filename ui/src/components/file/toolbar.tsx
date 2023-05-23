import { ChangeEvent, useCallback, useEffect, useRef, useState } from 'react'
import { useParams, useSearchParams } from 'react-router-dom'
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
} from '@koupr/ui'
import { BsSortDown } from 'react-icons/bs'
import FileAPI, { FileList } from '@/api/file'
import { ltEditorPermission, ltOwnerPermission } from '@/api/permission'
import downloadFile from '@/helpers/download-file'
import mapFileList from '@/helpers/map-file-list'
import { decodeQuery } from '@/helpers/query'
import { listUpdated } from '@/store/entities/files'
import { uploadAdded, UploadDecorator } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
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

const ICON_SCALE_LOCAL_STORAGE_KEY = 'voltaserve_file_icon_scale'
const SPACING = variables.spacingXs

const Toolbar = () => {
  const dispatch = useAppDispatch()
  const params = useParams()
  const workspaceId = params.id as string
  const fileId = params.fileId as string
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const selectionCount = useAppSelector(
    (state) => state.ui.files.selection.length
  )
  const singleFile = useAppSelector((state) =>
    state.ui.files.selection.length === 1
      ? state.entities.files.list?.data.find(
          (f) => f.id === state.ui.files.selection[0]
        )
      : null
  )
  const folder = useAppSelector((state) => state.entities.files.folder)
  const files = useAppSelector((state) => state.entities.files.list?.data)
  const iconScale = useAppSelector((state) => state.ui.files.iconScale)
  const hasOwnerPermission = useAppSelector(
    (state) =>
      state.entities.files.list?.data.findIndex(
        (f) =>
          state.ui.files.selection.findIndex(
            (s) => f.id === s && ltOwnerPermission(f.permission)
          ) !== -1
      ) === -1
  )
  const hasEditorPermission = useAppSelector(
    (state) =>
      state.entities.files.list?.data.findIndex(
        (f) =>
          state.ui.files.selection.findIndex(
            (s) => f.id === s && ltEditorPermission(f.permission)
          ) !== -1
      ) === -1
  )
  const uploadHiddenInput = useRef<HTMLInputElement>(null)

  useEffect(() => {
    const iconScale = localStorage.getItem(ICON_SCALE_LOCAL_STORAGE_KEY)
    if (iconScale) {
      dispatch(iconScaleUpdated(JSON.parse(iconScale)))
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
            }).value
          )
        )
      }
      dispatch(uploadsDrawerOpened())
      if (uploadHiddenInput && uploadHiddenInput.current) {
        uploadHiddenInput.current.value = ''
      }
    },
    [workspaceId, fileId, dispatch]
  )

  const handleIconScaleChange = useCallback(
    (value: number) => {
      localStorage.setItem(ICON_SCALE_LOCAL_STORAGE_KEY, JSON.stringify(value))
      dispatch(iconScaleUpdated(value))
    },
    [dispatch]
  )

  const handleRefresh = useCallback(async () => {
    setIsRefreshing(true)
    dispatch(selectionUpdated([]))
    try {
      let result: FileList
      if (query) {
        result = await FileAPI.search(
          { text: query, parentId: fileId, workspaceId },
          FileAPI.DEFAULT_PAGE_SIZE,
          1
        )
      } else {
        result = await FileAPI.list(fileId, FileAPI.DEFAULT_PAGE_SIZE, 1)
      }
      dispatch(listUpdated(result))
    } finally {
      setIsRefreshing(false)
    }
  }, [dispatch, fileId, workspaceId, query])

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
                    Share
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
        <Stack direction="row" spacing={SPACING}>
          <Slider
            w="120px"
            value={iconScale}
            min={1}
            max={2.5}
            step={0.25}
            onChange={handleIconScaleChange}
          >
            <SliderTrack>
              <Box position="relative" />
              <SliderFilledTrack />
            </SliderTrack>
            <SliderThumb />
          </Slider>
          <IconButton
            icon={<BsSortDown />}
            fontSize="16px"
            variant="solid"
            aria-label=""
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
                  <MenuItem icon={<IconCheckCircle />}>Sort By Name</MenuItem>
                  <MenuItem icon={<IconCircle />}>
                    Sort By Date Created
                  </MenuItem>
                  <MenuItem icon={<IconCircle />}>
                    Sort By Date Modified
                  </MenuItem>
                </MenuList>
              </Portal>
            </Menu>
          </Box>
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
