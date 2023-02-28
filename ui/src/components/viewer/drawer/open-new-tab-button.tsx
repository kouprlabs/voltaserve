import { useMemo } from 'react'
import { Button, IconButton } from '@chakra-ui/react'
import { IconExternalLink, variables } from '@koupr/ui'
import { File } from '@/api/file'

type OpenNewTabButtonProps = {
  file: File
  isCollapsed?: boolean
}

const LABEL = 'Open file'

const OpenNewTabButton = ({ file, isCollapsed }: OpenNewTabButtonProps) => {
  const url = useMemo(() => {
    if (file.original?.extension) {
      return `/proxy/api/v1/files/${file.id}/original${file.original.extension}`
    } else {
      return ''
    }
  }, [file])
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

export default OpenNewTabButton
