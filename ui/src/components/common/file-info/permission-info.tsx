import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'

type PermissionInfoProps = {
  file: File
}

const PermissionInfo = ({ file }: PermissionInfoProps) => (
  <Stat>
    <StatLabel>Permission</StatLabel>
    <StatNumber fontSize={variables.bodyFontSize}>
      <Badge>{file.permission}</Badge>
    </StatNumber>
  </Stat>
)

export default PermissionInfo
