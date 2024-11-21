// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Button, IconButton, Tooltip } from '@chakra-ui/react'
import { IconDownload } from '@koupr/ui'
import cx from 'classnames'
import { File } from '@/client/api/file'
import downloadFile from '@/lib/helpers/download-file'

export type DrawerDownloadButtonProps = {
  file: File
  isCollapsed?: boolean
}

const DrawerDownloadButton = ({
  file,
  isCollapsed,
}: DrawerDownloadButtonProps) => {
  const label = 'Download'
  return (
    <Tooltip label={label} isDisabled={!isCollapsed}>
      {isCollapsed ? (
        <IconButton
          icon={<IconDownload />}
          variant="solid"
          colorScheme="blue"
          title="Download"
          aria-label="Download"
          className={cx('h-[50px]', 'w-[50px]', 'p-1.5', 'rounded-md')}
          onClick={() => downloadFile(file)}
        />
      ) : (
        <Button
          leftIcon={<IconDownload />}
          variant="solid"
          colorScheme="blue"
          className={cx('h-[50px]', 'w-full', 'p-1.5', 'rounded-md')}
          onClick={() => downloadFile(file)}
        >
          {label}
        </Button>
      )}
    </Tooltip>
  )
}

export default DrawerDownloadButton
