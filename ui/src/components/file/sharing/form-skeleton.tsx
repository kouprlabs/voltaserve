import { Skeleton } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import classNames from 'classnames'

const FormSkeleton = () => (
  <div className={classNames('flex', 'flex-col', 'gap-1.5')}>
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
  </div>
)

export default FormSkeleton
