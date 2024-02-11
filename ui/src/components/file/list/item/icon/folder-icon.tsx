import { useMemo } from 'react'
import { Box, HStack } from '@chakra-ui/react'
import { FcFolder } from 'react-icons/fc'
import { CommonItemProps } from '@/types/file'
import ProcessingBadge from './processing-badge'
import SharedBadge from './shared-badge'

type FolderIconProps = {
  isLoading?: boolean
} & CommonItemProps

const ICON_FONT_SIZE = 92

const FolderIcon = ({ file, scale, viewType, isLoading }: FolderIconProps) => {
  const fontSize = useMemo(() => `${ICON_FONT_SIZE * scale}px`, [scale])
  const { bottom, right } = useMemo(() => {
    if (viewType === 'grid') {
      return { bottom: '7px', right: '2px' }
    } else {
      return { bottom: '0px', right: '-2px' }
    }
  }, [viewType])
  return (
    <Box position="relative">
      <FcFolder fontSize={fontSize} />
      <HStack position="absolute" bottom={bottom} right={right} spacing="2px">
        {file.isShared ? <SharedBadge /> : null}
        {isLoading ? <ProcessingBadge /> : null}
      </HStack>
    </Box>
  )
}

export default FolderIcon
