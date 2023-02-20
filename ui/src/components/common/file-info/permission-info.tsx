import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { File } from '@/api/file'
import variables from '@/theme/variables'

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
