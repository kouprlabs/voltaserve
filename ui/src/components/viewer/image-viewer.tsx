import { useMemo, useState } from 'react'
import { Center, Stack } from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import { File } from '@/api/file'

type ImageViewerProps = {
  file: File
}

const ImageViewer = ({ file }: ImageViewerProps) => {
  const [isLoading, setIsLoading] = useState(true)
  const download = useMemo(() => file.preview || file.original, [file])
  const urlPath = useMemo(() => (file.preview ? 'preview' : 'original'), [file])
  const url = useMemo(() => {
    if (!download || !download.extension) {
      return ''
    }
    return `/proxy/api/v1/files/${file.id}/${urlPath}${download.extension}`
  }, [file, download, urlPath])

  if (!download) {
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
