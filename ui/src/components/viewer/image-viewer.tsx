import { useMemo, useState } from 'react'
import { Center, Stack } from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import { File } from '@/api/file'

type ImageViewerProps = {
  file: File
}

const ImageViewer = ({ file }: ImageViewerProps) => {
  const [isLoading, setIsLoading] = useState(true)
  const url = useMemo(
    () => `/proxy/api/v1/files/${file.id}/original${file!.original!.extension}`,
    [file]
  )
  if (!file.original?.image) {
    return null
  }
  return (
    <Stack direction="column" w="100%" h="100%" spacing={variables.spacing}>
      <Center
        flexGrow={1}
        w="100%"
        h="100%"
        overflow="scroll"
        position="relative"
      >
        {isLoading && <SectionSpinner />}
        <img
          src={url}
          width={file.original.image.width}
          height={file.original.image.height}
          style={{
            objectFit: 'contain',
            width: isLoading ? 0 : 'auto',
            height: isLoading ? 0 : '100%',
            visibility: isLoading ? 'hidden' : 'visible',
          }}
          onLoad={() => setIsLoading(false)}
          alt={file.name}
        />
      </Center>
    </Stack>
  )
}

export default ImageViewer
