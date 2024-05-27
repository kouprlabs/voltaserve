import { useEffect, useMemo, useRef, useState } from 'react'
import { useColorMode } from '@chakra-ui/system'
import { Select } from 'chakra-react-select'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { Metadata } from '@/client/api/mosaic'
import { getConfig } from '@/config/config'
import { getAccessTokenOrRedirect } from '@/infra/token'
import reactSelectStyles from '@/styles/react-select'

export type ViewerImageProps = {
  file: File
}

const ViewerMosaic = ({ file }: ViewerImageProps) => {
  const accessToken = useMemo(() => getAccessTokenOrRedirect(), [])
  const { colorMode } = useColorMode()
  const [metadata, setMetadata] = useState<Metadata | null>(null)
  const [zoomLevel, setZoomLevel] = useState<number>(0)
  const [dragging, setDragging] = useState<boolean>(false)
  const [offset, setOffset] = useState<{ x: number; y: number }>({ x: 0, y: 0 })
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const offscreenCanvasRef = useRef<HTMLCanvasElement>(
    document.createElement('canvas'),
  )

  useEffect(() => {
    const fetchMetadata = async () => {
      const response = await fetch(
        `${getConfig().apiURL}/mosaics/${file.id}/metadata`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      )
      const data: Metadata = await response.json()
      setMetadata(data)
    }
    fetchMetadata()
  }, [file, accessToken])

  useEffect(() => {
    const renderCanvas = async () => {
      if (!metadata || !canvasRef.current || !offscreenCanvasRef.current) return

      const canvas = canvasRef.current
      const ctx = canvas.getContext('2d')
      const offscreenCanvas = offscreenCanvasRef.current
      const offscreenCtx = offscreenCanvas.getContext('2d')

      if (!ctx || !offscreenCtx) return

      const currentZoomLevel = metadata.zoomLevels[zoomLevel]

      offscreenCanvas.width = currentZoomLevel.width
      offscreenCanvas.height = currentZoomLevel.height

      // Clear the offscreen canvas
      offscreenCtx.clearRect(
        0,
        0,
        offscreenCanvas.width,
        offscreenCanvas.height,
      )

      const tilePromises = []

      for (let row = 0; row < currentZoomLevel.rows; row++) {
        for (let col = 0; col < currentZoomLevel.cols; col++) {
          const tileWidth =
            col === currentZoomLevel.cols - 1
              ? currentZoomLevel.tile.lastColWidth
              : currentZoomLevel.tile.width
          const tileHeight =
            row === currentZoomLevel.rows - 1
              ? currentZoomLevel.tile.lastRowHeight
              : currentZoomLevel.tile.height
          const img = new Image()
          const extension =
            file.snapshot?.preview?.extension ||
            file.snapshot?.original?.extension
          img.src = `${getConfig().apiURL}/mosaics/${file.id}/zoom_level/${zoomLevel}/row/${row}/col/${col}/ext/${extension?.replaceAll('.', '')}?access_token=${accessToken}`

          const promise = new Promise<void>((resolve) => {
            img.onload = () => {
              offscreenCtx.drawImage(
                img,
                col * currentZoomLevel.tile.width,
                row * currentZoomLevel.tile.height,
                tileWidth,
                tileHeight,
              )
              resolve()
            }
          })

          tilePromises.push(promise)
        }
      }

      await Promise.all(tilePromises)

      /* Clear the visible canvas */
      ctx.clearRect(0, 0, canvas.width, canvas.height)
      ctx.drawImage(offscreenCanvas, offset.x, offset.y)
    }

    renderCanvas()
  }, [metadata, zoomLevel, offset, file, accessToken])

  const handleMouseDown = () => {
    setDragging(true)
  }

  const handleMouseUp = () => {
    setDragging(false)
  }

  const handleMouseMove = (e: React.MouseEvent) => {
    if (dragging) {
      setOffset((prev) => ({
        x: prev.x + e.movementX,
        y: prev.y + e.movementY,
      }))
    }
  }

  const handleZoomChange = (value: number) => {
    setZoomLevel(value)
    setOffset({ x: 0, y: 0 })
  }

  if (!metadata) {
    return null
  }

  return (
    <div className={cx('absolute', 'top-0', 'left-0')}>
      <Select
        className={cx(
          'absolute',
          'top-[15px]',
          'left-[15px]',
          'z-10',
          'w-[200px]',
        )}
        defaultValue={{
          value: metadata.zoomLevels[0].index,
          label: `Zoom ${metadata.zoomLevels[0].scaleDownPercentage}%`,
        }}
        options={metadata.zoomLevels.map((zoomLevel) => ({
          value: zoomLevel.index,
          label: `Zoom ${zoomLevel.scaleDownPercentage}%`,
        }))}
        placeholder="Zoom Level"
        selectedOptionStyle="check"
        chakraStyles={reactSelectStyles({ colorMode })}
        onChange={(event) => {
          if (event?.value !== undefined) {
            handleZoomChange(event.value)
          }
        }}
      />
      <canvas
        className={cx('absolute', 'top-0', 'left-0', 'z-0')}
        ref={canvasRef}
        width={window.innerWidth}
        height={window.innerHeight}
        onMouseDown={handleMouseDown}
        onMouseUp={handleMouseUp}
        onMouseMove={handleMouseMove}
      ></canvas>
    </div>
  )
}

export default ViewerMosaic
