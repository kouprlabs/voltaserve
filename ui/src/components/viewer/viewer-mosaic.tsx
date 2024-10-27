// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useMemo, useRef, useState, MouseEvent } from 'react'
import { useColorMode } from '@chakra-ui/system'
import { Select } from 'chakra-react-select'
import cx from 'classnames'
import { File } from '@/client/api/file'
import MosaicAPI, { Metadata, ZoomLevel } from '@/client/api/mosaic'
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
  const tileCache = useRef<Map<string, HTMLImageElement>>(new Map())

  useEffect(() => {
    ;(async function (file: File) {
      const { metadata } = await MosaicAPI.getInfo(file.id)
      if (metadata) {
        setMetadata(metadata)
      }
    })(file)
  }, [file])

  useEffect(() => {
    ;(async function () {
      if (!metadata || !canvasRef.current) return

      const canvas = canvasRef.current
      const ctx = canvas.getContext('2d')

      if (!ctx) return

      const currentZoomLevel = metadata.zoomLevels[zoomLevel]
      canvas.width = window.innerWidth
      canvas.height = window.innerHeight

      // Clear the visible canvas
      ctx.clearRect(0, 0, canvas.width, canvas.height)

      const visibleTiles = getVisibleTiles(
        currentZoomLevel,
        canvas.width,
        canvas.height,
        offset,
      )

      visibleTiles.forEach(({ row, col }) => {
        const tileKey = `${zoomLevel}-${row}-${col}`
        const cachedTile = tileCache.current.get(tileKey)

        if (cachedTile) {
          drawTile(ctx, cachedTile, currentZoomLevel, row, col)
        } else {
          const img = new Image()
          let extension =
            file.snapshot?.preview?.extension ||
            file.snapshot?.original?.extension
          extension = extension?.replaceAll('.', '')
          img.src = `${getConfig().apiURL}/mosaics/${file.id}/zoom_level/${zoomLevel}/row/${row}/col/${col}/ext/${extension}?access_token=${accessToken}`

          img.onload = () => {
            tileCache.current.set(tileKey, img)
            drawTile(ctx, img, currentZoomLevel, row, col)
          }
        }
      })
    })()
  }, [metadata, zoomLevel, offset, file, accessToken])

  const drawTile = (
    ctx: CanvasRenderingContext2D,
    img: HTMLImageElement,
    zoomLevel: ZoomLevel,
    row: number,
    col: number,
  ) => {
    const tileWidth =
      col === zoomLevel.cols - 1
        ? zoomLevel.tile.lastColWidth
        : zoomLevel.tile.width
    const tileHeight =
      row === zoomLevel.rows - 1
        ? zoomLevel.tile.lastRowHeight
        : zoomLevel.tile.height
    ctx.drawImage(
      img,
      col * zoomLevel.tile.width + offset.x,
      row * zoomLevel.tile.height + offset.y,
      tileWidth,
      tileHeight,
    )
  }

  const getVisibleTiles = (
    zoomLevel: ZoomLevel,
    viewportWidth: number,
    viewportHeight: number,
    offset: { x: number; y: number },
  ) => {
    const tiles: { row: number; col: number }[] = []
    const startX = Math.max(0, Math.floor(-offset.x / zoomLevel.tile.width))
    const startY = Math.max(0, Math.floor(-offset.y / zoomLevel.tile.height))
    const endX = Math.min(
      zoomLevel.cols,
      Math.ceil((viewportWidth - offset.x) / zoomLevel.tile.width),
    )
    const endY = Math.min(
      zoomLevel.rows,
      Math.ceil((viewportHeight - offset.y) / zoomLevel.tile.height),
    )

    for (let row = startY; row < endY; row++) {
      for (let col = startX; col < endX; col++) {
        tiles.push({ row, col })
      }
    }
    return tiles
  }

  const handleMouseDown = () => {
    setDragging(true)
  }

  const handleMouseUp = () => {
    setDragging(false)
  }

  const handleMouseMove = (e: MouseEvent) => {
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
          label: `Zoom ${Math.round(metadata.zoomLevels[0].scaleDownPercentage)}%`,
        }}
        options={metadata.zoomLevels.map((zoomLevel) => ({
          value: zoomLevel.index,
          label: `Zoom ${Math.round(zoomLevel.scaleDownPercentage)}%`,
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
