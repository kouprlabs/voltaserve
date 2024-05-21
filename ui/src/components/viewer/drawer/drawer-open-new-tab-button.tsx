import { useMemo } from 'react'
import { Button, IconButton } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { IconOpenInNew } from '@/lib'

export type DrawerOpenNewTabButtonProps = {
  file: File
  isCollapsed?: boolean
}

const DrawerOpenNewTabButton = ({
  file,
  isCollapsed,
}: DrawerOpenNewTabButtonProps) => {
  const label = 'Open file'
  const download = useMemo(() => file.preview ?? file.original, [file])
  const path = useMemo(() => (file.preview ? 'preview' : 'original'), [file])
  const url = useMemo(() => {
    if (!download?.extension) {
      return ''
    }
    if (file.original?.extension) {
      return `/proxy/api/v2/files/${file.id}/${path}${download.extension}`
    } else {
      return ''
    }
  }, [file, download, path])
  if (!file.original) {
    return null
  }
  if (isCollapsed) {
    return (
      <IconButton
        icon={<IconOpenInNew />}
        as="a"
        className={cx('h-[50px]', 'w-[50px]', 'p-1.5', 'rounded-md')}
        href={url}
        target="_blank"
        title={label}
        aria-label={label}
      />
    )
  } else {
    return (
      <Button
        leftIcon={<IconOpenInNew />}
        as="a"
        className={cx('h-[50px]', 'w-full', 'p-1.5', 'rounded-md')}
        href={url}
        target="_blank"
      >
        {label}
      </Button>
    )
  }
}

export default DrawerOpenNewTabButton
