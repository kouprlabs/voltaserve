import { Avatar } from '@chakra-ui/react'
import cx from 'classnames'
import { randomColorFromString } from '@/lib/helpers/colorGenerator'

export interface UserAvatarProps {
  name: string
  src?: string
  height?: string
  width?: string
  className?: string
  size?: string
}

const UserAvatar = ({
  name,
  src,
  size = 'sm',
  className,
  width,
  height = '40px',
}: UserAvatarProps) => {
  return src !== '' ? (
    <Avatar
      name={name}
      src={src}
      size={size}
      className={cx(
        `w-[${width ? width : height}]`,
        `h-[${height}]`,
        'border',
        'border-gray-300',
        'dark:border-gray-700',
        className,
      )}
    />
  ) : (
    <Avatar
      name={name}
      size={size}
      className={cx(
        `w-[${width ? width : height}]`,
        `h-[${height}]`,
        'border',
        'border-gray-300',
        'dark:border-gray-700',
        className,
      )}
      sx={{ backgroundColor: randomColorFromString(name) }}
    />
  )
}

export default UserAvatar
