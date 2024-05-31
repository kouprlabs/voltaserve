import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { IconSecurity } from '@/lib/components/icons'

const IconBadgeSecurity = () => (
  <Tooltip label="This item has a watermark">
    <Circle
      className={cx(
        'text-orange-600',
        'bg-white',
        'w-[23px]',
        'h-[23px]',
        'border',
        'border-gray-200',
      )}
    >
      <IconSecurity className={cx('text-[12px]')} />
    </Circle>
  </Tooltip>
)

export default IconBadgeSecurity
