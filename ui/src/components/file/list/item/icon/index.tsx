import { useColorModeValue, useToken } from '@chakra-ui/react'
import classNames from 'classnames'
import { CommonItemProps } from '@/types/file'
import FileIcon from './file-icon'
import FolderIcon from './folder-icon'

type IconProps = {
  isLoading?: boolean
} & CommonItemProps

const Icon = ({ file, scale, viewType, isLoading }: IconProps) => {
  const color = useToken('colors', useColorModeValue('gray.500', 'gray.300'))
  return (
    <>
      {file.type === 'file' ? (
        <div className={classNames('z-0')} style={{ color }}>
          <FileIcon file={file} scale={scale} viewType={viewType} />
        </div>
      ) : file.type === 'folder' ? (
        <div className={classNames('z-0')} style={{ color }}>
          <FolderIcon
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

export default Icon
