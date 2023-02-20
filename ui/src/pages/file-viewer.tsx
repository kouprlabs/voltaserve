import { useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { Button, Center, Heading, Spinner, Stack, Text } from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import FileAPI, { File } from '@/api/file'
import Drawer from '@/components/common/drawer'
import { IconDownload } from '@/components/common/icon'
import AudioPlayer from '@/components/viewer/audio-player'
import DrawerContent from '@/components/viewer/drawer/content'
import ImageViewer from '@/components/viewer/image-viewer'
import PdfViewer from '@/components/viewer/pdf-viewer'
import VideoPlayer from '@/components/viewer/video-player'
import variables from '@/theme/variables'
import downloadFile from '@/helpers/download-file'
import { isAudio, isImage, isPdf, isVideo } from '@/helpers/file-extension'

const FileViewerPage = () => {
  const params = useParams()
  const { data: file } = FileAPI.useGetById(params.id as string)

  const renderViewer = useCallback((file: File) => {
    if (file.preview && isPdf(file.preview.extension)) {
      return <PdfViewer file={file} />
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
        <Spinner size="sm" thickness="4px" />
      </Center>
    )
  }

  return (
    <>
      <Helmet>
        <title>{file.name}</title>
      </Helmet>
      <Stack direction="row" spacing={0} h="100%">
        <Drawer localStorageNamespace="viewer">
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
