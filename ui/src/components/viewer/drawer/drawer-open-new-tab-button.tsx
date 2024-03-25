import { useMemo } from 'react'
import { Button, IconButton } from '@chakra-ui/react'
import { IconExternalLink, variables } from '@koupr/ui'
import { File } from '@/client/api/file'

type DrawerOpenNewTabButtonProps = {
  file: File
  isCollapsed?: boolean
}

const LABEL = 'Open file'

const DrawerOpenNewTabButton = ({
  file,
  isCollapsed,
}: DrawerOpenNewTabButtonProps) => {
  const download = useMemo(() => file.preview ?? file.original, [file])
  const path = useMemo(() => (file.preview ? 'preview' : 'original'), [file])
  const url = useMemo(() => {
    if (!download?.extension) {
      return ''
    }
    if (file.original?.extension) {
      return `/proxy/api/v1/files/${file.id}/${path}${download.extension}`
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
        icon={<IconExternalLink fontSize="18px" />}
        as="a"
        w="50px"
        h="50px"
        p={variables.spacing}
        borderRadius={variables.borderRadiusSm}
        href={url}
        target="_blank"
        title={LABEL}
        aria-label={LABEL}
      />
    )
  } else {
    return (
      <Button
        leftIcon={<IconExternalLink fontSize="18px" />}
        as="a"
        w="100%"
        h="50px"
        p={variables.spacing}
        borderRadius={variables.borderRadiusSm}
        href={url}
        target="_blank"
      >
        {LABEL}
      </Button>
    )
  }
}

export default DrawerOpenNewTabButton
