// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useMemo } from 'react'
import { Button, IconButton, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { IconOpenInNew } from '@/lib/components/icons'

export type DrawerOpenNewTabButtonProps = {
  file: File
  isCollapsed?: boolean
}

const DrawerOpenNewTabButton = ({
  file,
  isCollapsed,
}: DrawerOpenNewTabButtonProps) => {
  const label = 'Open file'
  const download = useMemo(
    () => file.snapshot?.preview ?? file.snapshot?.original,
    [file],
  )
  const path = useMemo(
    () => (file.snapshot?.preview ? 'preview' : 'original'),
    [file],
  )
  const url = useMemo(() => {
    if (!download?.extension) {
      return ''
    }
    if (file.snapshot?.original?.extension) {
      return `/proxy/api/v3/files/${file.id}/${path}${download.extension}`
    } else {
      return ''
    }
  }, [file, download, path])
  if (!file.snapshot?.original) {
    return null
  }
  return (
    <Tooltip label={label} isDisabled={!isCollapsed}>
      {isCollapsed ? (
        <IconButton
          icon={<IconOpenInNew />}
          as="a"
          className={cx('h-[50px]', 'w-[50px]', 'p-1.5', 'rounded-md')}
          href={url}
          target="_blank"
          title={label}
          aria-label={label}
        />
      ) : (
        <Button
          leftIcon={<IconOpenInNew />}
          as="a"
          className={cx('h-[50px]', 'w-full', 'p-1.5', 'rounded-md')}
          href={url}
          target="_blank"
        >
          {label}
        </Button>
      )}
    </Tooltip>
  )
}

export default DrawerOpenNewTabButton
