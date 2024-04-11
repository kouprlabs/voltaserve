import { Skeleton } from '@chakra-ui/react'
import cx from 'classnames'

const SharingFormSkeleton = () => (
  <div className={cx('flex', 'flex-col', 'gap-1.5')}>
    <Skeleton height="40px" className={cx('rounded-xl')} />
    <Skeleton height="40px" className={cx('rounded-xl')} />
    <Skeleton height="40px" className={cx('rounded-xl')} />
  </div>
)

export default SharingFormSkeleton
