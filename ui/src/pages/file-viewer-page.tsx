import { useCallback, useMemo } from 'react'
import { useLocation, useParams } from 'react-router-dom'
import { Button } from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import FileAPI, { File } from '@/client/api/file'
import DrawerContent from '@/components/viewer/drawer/drawer-content'
import ViewerAudio from '@/components/viewer/viewer-audio'
import ViewerImage from '@/components/viewer/viewer-image'
import ViewerMosaic from '@/components/viewer/viewer-mosaic'
import ViewerPDF from '@/components/viewer/viewer-pdf'
import ViewerVideo from '@/components/viewer/viewer-video'
import downloadFile from '@/helpers/download-file'
import { isAudio, isImage, isPDF, isVideo } from '@/helpers/file-extension'
import Drawer from '@/lib/components/drawer'
import { IconDownload } from '@/lib/components/icons'
import Spinner from '@/lib/components/spinner'

const FileViewerPage = () => {
  const { id } = useParams()
  const location = useLocation()
  const isMosaic = useMemo(
    () => location.pathname.endsWith('/mosaic'),
    [location],
  )
  const { data: file } = FileAPI.useGet(id)

  const renderViewer = useCallback(
    (file: File) => {
      if (isMosaic) {
        return <ViewerMosaic file={file} />
      } else {
        if (
          (file.snapshot?.original &&
            isPDF(file.snapshot?.original.extension)) ||
          (file.snapshot?.preview && isPDF(file.snapshot?.preview.extension))
        ) {
          return <ViewerPDF file={file} />
        } else if (
          file.snapshot?.original &&
          isImage(file.snapshot?.original.extension)
        ) {
          if (file.snapshot?.mosaic) {
            return <ViewerMosaic file={file} />
          } else {
            return <ViewerImage file={file} />
          }
        } else if (
          file.snapshot?.original &&
          isVideo(file.snapshot?.original.extension)
        ) {
          return <ViewerVideo file={file} />
        } else if (
          file.snapshot?.original &&
          isAudio(file.snapshot?.original.extension)
        ) {
          return <ViewerAudio file={file} />
        } else {
          return (
            <div className={cx('flex', 'flex-col', 'gap-1.5')}>
              <span className={cx('text-[16px]')}>
                Cannot preview this file.
              </span>
              <Button
                leftIcon={<IconDownload />}
                colorScheme="blue"
                onClick={() => downloadFile(file)}
              >
                Download
              </Button>
            </div>
          )
        }
      }
    },
    [isMosaic],
  )

  return (
    <>
      {file ? (
        <>
          <Helmet>
            <title>{file.name}</title>
          </Helmet>
          <div className={cx('flex', 'flex-row', 'gap-0', 'h-full')}>
            <Drawer storage={{ prefix: 'voltaserve', namespace: 'viewer' }}>
              <DrawerContent file={file} />
            </Drawer>
            <div
              className={cx('flex', 'flex-col', 'gap-0', 'grow', 'h-[100vh]')}
            >
              <div
                className={cx(
                  'flex',
                  'items-center',
                  'justify-center',
                  'w-full',
                  'h-[80px]',
                  'min-h-[80px]',
                )}
              >
                <span className={cx('font-medium', 'text-[16px]')}>
                  {file.name}
                </span>
              </div>
              <div
                className={cx(
                  'flex',
                  'items-center',
                  'justify-center',
                  'w-full',
                  'h-full',
                  'overflow-hidden',
                  'relative',
                )}
              >
                {renderViewer(file)}
              </div>
            </div>
          </div>
        </>
      ) : (
        <div
          className={cx('flex', 'items-center', 'justify-center', 'h-[100vh]')}
        >
          <Spinner />
        </div>
      )}
    </>
  )
}

export default FileViewerPage
