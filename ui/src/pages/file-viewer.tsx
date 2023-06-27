import { useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { Button, Center, Heading, Stack, Text } from '@chakra-ui/react'
import { variables, IconDownload, Drawer, Spinner } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import FileAPI, { File } from '@/client/api/file'
import AudioPlayer from '@/components/viewer/audio-player'
import DrawerContent from '@/components/viewer/drawer/content'
import ImageViewer from '@/components/viewer/image-viewer'
import OCRViewer from '@/components/viewer/ocr-viewer'
import PDFViewer from '@/components/viewer/pdf-viewer'
import VideoPlayer from '@/components/viewer/video-player'
import downloadFile from '@/helpers/download-file'
import { isAudio, isImage, isPDF, isVideo } from '@/helpers/file-extension'

const FileViewerPage = () => {
  const params = useParams()
  const { data: file } = FileAPI.useGetById(params.id as string)

  const renderViewer = useCallback((file: File) => {
    if (
      (file.original && isPDF(file.original.extension)) ||
      (file.preview && isPDF(file.preview.extension))
    ) {
      return <PDFViewer file={file} />
    } else if (file.original && isImage(file.original.extension) && file.ocr) {
      return <OCRViewer file={file} />
    } else if (file.original && isImage(file.original.extension)) {
      return <ImageViewer file={file} />
    } else if (file.original && isVideo(file.original.extension)) {
      return <VideoPlayer file={file} />
    } else if (file.original && isAudio(file.original.extension)) {
      return <AudioPlayer file={file} />
    } else {
      return (
        <Stack direction="column" spacing={variables.spacing}>
          <Text fontSize="16px">Cannot preview this file.</Text>
          <Button
            leftIcon={<IconDownload />}
            colorScheme="blue"
            onClick={() => downloadFile(file)}
          >
            Download
          </Button>
        </Stack>
      )
    }
  }, [])

  if (!file) {
    return (
      <Center height="100vh">
        <Spinner />
      </Center>
    )
  }

  return (
    <>
      <Helmet>
        <title>{file.name}</title>
      </Helmet>
      <Stack direction="row" spacing={0} h="100%">
        <Drawer localStoragePrefix="voltaserve" localStorageNamespace="viewer">
          <DrawerContent file={file} />
        </Drawer>
        <Stack height="100vh" spacing={0} flexGrow={1}>
          <Center w="100%" h="80px">
            <Heading fontSize="14px" textAlign="center">
              {file.name}
            </Heading>
          </Center>
          <Center w="100%" h="100%" overflow="hidden">
            {renderViewer(file)}
          </Center>
        </Stack>
      </Stack>
    </>
  )
}

export default FileViewerPage
