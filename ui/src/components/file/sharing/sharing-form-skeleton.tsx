import { Skeleton } from '@chakra-ui/react'
import cx from 'classnames'

const SharingFormSkeleton = () => (
  <div className={cx('flex', 'flex-col', 'gap-1.5')}>
    <Skeleton className={cx('rounded-xl', 'w-[40px]')} />
    <Skeleton className={cx('rounded-xl', 'w-[40px]')} />
    <Skeleton className={cx('rounded-xl', 'w-[40px]')} />
  </div>
)

export default SharingFormSkeleton
