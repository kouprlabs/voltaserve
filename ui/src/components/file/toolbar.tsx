import { ChangeEvent, useCallback, useEffect, useRef } from 'react'
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
import { variables } from '@koupr/ui'
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
} from '@koupr/ui'
import { ltEditorPermission, ltOwnerPermission } from '@/api/permission'
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
import downloadFile from '@/helpers/download-file'
import mapFileList from '@/helpers/map-file-list'

const ICON_SCALE_LOCAL_STORAGE_KEY = 'voltaserve_file_icon_scale'

const Toolbar = () => {
  const dispatch = useAppDispatch()
  const params = useParams()
  const workspaceId = params.id as string
  const fileId = params.fileId as string
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

  return (
    <>
      <Stack direction="row" spacing={variables.spacingSm}>
        <ButtonGroup isAttached>
          <Button
            variant="solid"
            colorScheme="blue"
            leftIcon={<IconUpload />}
            isDisabled={!folder || ltEditorPermission(folder.permission)}
            onClick={() => uploadHiddenInput?.current?.click()}
          >
            Upload file
          </Button>
          <Button
            variant="outline"
            colorScheme="blue"
            leftIcon={<IconAdd />}
            isDisabled={!folder || ltEditorPermission(folder.permission)}
            onClick={() => dispatch(createModalDidOpen())}
          >
            New folder
          </Button>
        </ButtonGroup>
        <Stack direction="row" spacing={variables.spacingSm}>
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
                    Select all
                  </MenuItem>
                  <MenuItem
                    icon={<IconCircle />}
                    onClick={() => dispatch(selectionUpdated([]))}
                  >
                    Unselect all
                  </MenuItem>
                </MenuList>
              </Portal>
            </Menu>
          </Box>
        </Stack>
        <Spacer />
        <Slider
          w="200px"
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
