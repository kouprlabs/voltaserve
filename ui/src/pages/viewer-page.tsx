import { useCallback, useMemo } from 'react'
import { useLocation, useParams } from 'react-router-dom'
import { Button } from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import FileAPI, { File } from '@/client/api/file'
import DrawerContent from '@/components/viewer/drawer/drawer-content'
import ViewerAudio from '@/components/viewer/viewer-audio'
import ViewerImage from '@/components/viewer/viewer-image'
import ViewerModel from '@/components/viewer/viewer-model'
import ViewerMosaic from '@/components/viewer/viewer-mosaic'
import ViewerPDF from '@/components/viewer/viewer-pdf'
import ViewerVideo from '@/components/viewer/viewer-video'
import Drawer from '@/lib/components/drawer'
import { IconDownload } from '@/lib/components/icons'
import Spinner from '@/lib/components/spinner'
import downloadFile from '@/lib/helpers/download-file'
import {
  isGLB,
  isAudio,
  isImage,
  isPDF,
  isVideo,
} from '@/lib/helpers/file-extension'

const ViewerPage = () => {
  const { id } = useParams()
  const location = useLocation()
  const { data: file } = FileAPI.useGet(id)
  const hasMosaicPath = useMemo(
    () => location.pathname.endsWith('/mosaic'),
    [location],
  )
  const hasPDF = useMemo(() => {
    return file?.snapshot &&
      ((file.snapshot.original && isPDF(file.snapshot.original.extension)) ||
        (file.snapshot.preview && isPDF(file.snapshot.preview.extension)))
      ? true
      : false
  }, [file, location])
  const hasImage = useMemo(
    () =>
      file?.snapshot &&
      file.snapshot?.original &&
      isImage(file.snapshot?.original.extension)
        ? true
        : false,
    [file],
  )
  const hasMosaicImage = useMemo(
    () => hasImage && file?.snapshot && file.snapshot?.mosaic,
    [hasImage],
  )
  const hasVideo = useMemo(
    () =>
      file?.snapshot &&
      file.snapshot?.original &&
      isVideo(file.snapshot?.original.extension),
    [file],
  )
  const hasAudio = useMemo(
    () =>
      file?.snapshot &&
      file.snapshot?.original &&
      isAudio(file.snapshot?.original.extension),
    [file],
  )
  const hasGLB = useMemo(
    () =>
      file?.snapshot &&
      file.snapshot?.preview &&
      isGLB(file.snapshot.preview.extension),
    [file],
  )

  const renderViewer = useCallback(
    (file: File) => {
      if (hasMosaicPath) {
        return <ViewerMosaic file={file} />
      } else {
        if (hasPDF) {
          return <ViewerPDF file={file} />
        } else if (hasImage) {
          if (hasMosaicImage) {
            return <ViewerMosaic file={file} />
          } else {
            return <ViewerImage file={file} />
          }
        } else if (hasVideo) {
          return <ViewerVideo file={file} />
        } else if (hasAudio) {
          return <ViewerAudio file={file} />
        } else if (hasGLB) {
          return <ViewerModel file={file} />
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
    [
      hasMosaicPath,
      hasMosaicImage,
      hasPDF,
      hasImage,
      hasVideo,
      hasAudio,
      hasGLB,
    ],
  )
  const isPresentational = hasVideo || (hasImage && !hasMosaicImage)

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
                  { 'bg-black': isPresentational },
                )}
              >
                <span
                  className={cx('font-medium', 'text-[16px]', {
                    'text-white': isPresentational,
                  })}
                >
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
                  { 'bg-black': isPresentational },
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

export default ViewerPage
