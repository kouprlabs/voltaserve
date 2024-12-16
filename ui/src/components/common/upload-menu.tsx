// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { Button, Menu, MenuButton, MenuItem, MenuList } from '@chakra-ui/react'
import { IconAdd, IconCloudUpload, IconUpload } from '@koupr/ui'
import FileAPI from '@/client/api/file'
import { geEditorPermission } from '@/client/api/permission'
import { useAppDispatch } from '@/store/hook'
import { fileUploadDidOpen, folderUploadDidOpen } from '@/store/ui/files'

const UploadMenu = () => {
  const dispatch = useAppDispatch()
  const { fileId } = useParams()
  const { data: folder } = FileAPI.useGet(fileId)

  const handleUploadFiles = useCallback(() => {
    dispatch(fileUploadDidOpen())
  }, [dispatch])

  const handleUploadFolders = useCallback(() => {
    dispatch(folderUploadDidOpen())
  }, [dispatch])

  return (
    <>
      {folder && geEditorPermission(folder.permission) ? (
        <Menu>
          <MenuButton
            as={Button}
            variant="solid"
            leftIcon={<IconCloudUpload />}
          >
            Upload
          </MenuButton>
          <MenuList>
            <MenuItem icon={<IconAdd />} onClick={handleUploadFiles}>
              Upload Files
            </MenuItem>
            <MenuItem icon={<IconUpload />} onClick={handleUploadFolders}>
              Upload Folder
            </MenuItem>
          </MenuList>
        </Menu>
      ) : null}
    </>
  )
}

export default UploadMenu
