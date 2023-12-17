import { Skeleton, Stack } from '@chakra-ui/react'
import { variables } from '@koupr/ui'

const FormSkeleton = () => (
  <Stack spacing={variables.spacing}>
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
  </Stack>
)

export default FormSkeleton
