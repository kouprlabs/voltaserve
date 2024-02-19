import { useToken } from '@chakra-ui/react'
import { DragOverlay } from '@dnd-kit/core'
import classNames from 'classnames'
import { useAppSelector } from '@/store/hook'
import { CommonItemProps } from '@/types/file'
import Item from './item'

type ItemDragOverlayProps = CommonItemProps

const ItemDragOverlay = ({ file, scale, viewType }: ItemDragOverlayProps) => {
  const selectionCount = useAppSelector(
    (state) => state.ui.files.selection.length,
  )
  const green = useToken('colors', 'green.300')

  return (
    <DragOverlay>
      <div className={classNames('relative')}>
        <Item
          file={file}
          scale={scale}
          isPresentational={true}
          viewType={viewType}
        />
        {selectionCount > 1 ? (
          <div
            className={classNames(
              'absolute',
              'flex',
              'items-center',
              'justify-center',
              'bottom-[-5px]',
              'right-[-5px]',
              'text-white',
              'rounded-xl',
              'min-w-[30px]',
              'h-[30px]',
              'px-1',
            )}
            style={{ backgroundColor: green }}
          >
            {selectionCount}
          </div>
        ) : null}
      </div>
    </DragOverlay>
  )
}

export default ItemDragOverlay
