import { useColorModeValue, useToken } from '@chakra-ui/react'
import cx from 'classnames'
import { FileCommonProps } from '@/types/file'
import IconFile from './icon-file'
import IconFolder from './icon-folder'

export type ItemIconProps = {
  isLoading?: boolean
} & FileCommonProps

const ItemIcon = ({ file, scale, viewType, isLoading }: ItemIconProps) => {
  const color = useToken('colors', useColorModeValue('gray.500', 'gray.300'))
  return (
    <>
      {file.type === 'file' ? (
        <div className={cx('z-0')} style={{ color }}>
          <IconFile file={file} scale={scale} viewType={viewType} />
        </div>
      ) : file.type === 'folder' ? (
        <div className={cx('z-0')} style={{ color }}>
          <IconFolder
            file={file}
            scale={scale}
            viewType={viewType}
            isLoading={isLoading}
          />
        </div>
      ) : null}
    </>
  )
}

export default ItemIcon
