// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect, useMemo, useRef, useState, MouseEvent, useCallback } from 'react'
import { Select } from '@koupr/ui'
import { OptionBase, SingleValue } from 'chakra-react-select'
import cx from 'classnames'
import { File } from '@/client/api/file'
import MosaicAPI, { Metadata, ZoomLevel } from '@/client/api/mosaic'
import { getConfig } from '@/config/config'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerImageProps = {
  file: File
}

interface ZoomLevelOption extends OptionBase {
  value: number
  label: string
}

type TileCoordinate = {
  row: number
  column: number
}

type MouseCoordinate = {
  x: number
  y: number
}

const ViewerMosaic = ({ file }: ViewerImageProps) => {
  const accessToken = useMemo(() => getAccessTokenOrRedirect(), [])
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
      if (!metadata) {
        return
      }
      const canvas = canvasRef.current
      if (!canvas) {
        return
      }
      const context = canvas.getContext('2d')
      if (!context) {
        return
      }
      const currentZoomLevel = metadata.zoomLevels[zoomLevel]
      context.clearRect(0, 0, canvas.width, canvas.height)
      getVisibleTileCoordinates(currentZoomLevel, canvas.width, canvas.height, offset).forEach(({ row, column }) => {
        const tileKey = `${zoomLevel}-${row}-${column}`
        const cachedTile = tileCache.current.get(tileKey)
        if (cachedTile) {
          drawTile(context, cachedTile, currentZoomLevel, row, column)
        } else {
          const image = new Image()
          let extension = file.snapshot?.preview?.extension || file.snapshot?.original?.extension
          extension = extension?.replaceAll('.', '')
          image.src = `${getConfig().apiURL}/mosaics/${file.id}/zoom_level/${zoomLevel}/row/${row}/column/${column}/extension/${extension}?access_token=${accessToken}`
          image.onload = () => {
            tileCache.current.set(tileKey, image)
            drawTile(context, image, currentZoomLevel, row, column)
          }
        }
      })
    })()
  }, [file, metadata, zoomLevel, offset, accessToken, canvasRef, tileCache])

  const drawTile = useCallback(
    (context: CanvasRenderingContext2D, image: HTMLImageElement, zoomLevel: ZoomLevel, row: number, column: number) => {
      const tileWidth = column === zoomLevel.cols - 1 ? zoomLevel.tile.lastColWidth : zoomLevel.tile.width
      const tileHeight = row === zoomLevel.rows - 1 ? zoomLevel.tile.lastRowHeight : zoomLevel.tile.height
      context.drawImage(
        image,
        column * zoomLevel.tile.width + offset.x,
        row * zoomLevel.tile.height + offset.y,
        tileWidth,
        tileHeight,
      )
    },
    [canvasRef, offset],
  )

  const getVisibleTileCoordinates = useCallback(
    (
      zoomLevel: ZoomLevel,
      viewportWidth: number,
      viewportHeight: number,
      offset: MouseCoordinate,
    ): TileCoordinate[] => {
      const coordinates: TileCoordinate[] = []
      const startX = Math.max(0, Math.floor(-offset.x / zoomLevel.tile.width))
      const startY = Math.max(0, Math.floor(-offset.y / zoomLevel.tile.height))
      const endX = Math.min(zoomLevel.cols, Math.ceil((viewportWidth - offset.x) / zoomLevel.tile.width))
      const endY = Math.min(zoomLevel.rows, Math.ceil((viewportHeight - offset.y) / zoomLevel.tile.height))
      for (let row = startY; row < endY; row++) {
        for (let column = startX; column < endX; column++) {
          coordinates.push({ row, column: column })
        }
      }
      return coordinates
    },
    [],
  )

  const handleMouseDown = useCallback(() => {
    setDragging(true)
  }, [])

  const handleMouseUp = useCallback(() => {
    setDragging(false)
  }, [])

  const handleMouseMove = useCallback(
    (event: MouseEvent) => {
      if (dragging) {
        setOffset((previous) => ({
          x: previous.x + event.movementX,
          y: previous.y + event.movementY,
        }))
      }
    },
    [dragging],
  )

  const handleZoomChange = useCallback((newValue: SingleValue<ZoomLevelOption>) => {
    if (newValue?.value !== undefined) {
      setZoomLevel(newValue.value)
      setOffset({ x: 0, y: 0 })
    }
  }, [])

  return (
    <>
      {metadata ? (
        <div className={cx('absolute', 'top-0', 'left-0')}>
          <Select<ZoomLevelOption, false>
            className={cx('absolute', 'top-1.5', 'left-1.5', 'z-10', 'w-[200px]')}
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
            onChange={handleZoomChange}
          />
          <canvas
            className={cx('absolute', 'top-0', 'left-0', 'z-0')}
            ref={canvasRef}
            width={window.innerWidth}
            height={window.innerHeight}
            onMouseDown={handleMouseDown}
            onMouseUp={handleMouseUp}
            onMouseMove={handleMouseMove}
          />
        </div>
      ) : null}
    </>
  )
}

export default ViewerMosaic
