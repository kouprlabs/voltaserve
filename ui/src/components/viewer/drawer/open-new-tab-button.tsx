import { useMemo } from 'react'
import { Button, IconButton } from '@chakra-ui/react'
import { File } from '@/api/file'
import { IconExternalLink } from '@/components/common/icon'
import variables from '@/theme/variables'

type OpenNewTabButtonProps = {
  file: File
  isCollapsed?: boolean
}

const LABEL = 'Open file'

const OpenNewTabButton = ({ file, isCollapsed }: OpenNewTabButtonProps) => {
  const url = useMemo(
    () => `/proxy/api/v1/files/${file.id}/original${file!.original!.extension}`,
    [file]
  )
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
