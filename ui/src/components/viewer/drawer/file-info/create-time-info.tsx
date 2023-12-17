import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'
import prettyDate from '@/helpers/pretty-date'

type CreateTimeInfoProps = {
  file: File
}

const CreateTimeInfo = ({ file }: CreateTimeInfoProps) => (
  <Stat>
    <StatLabel>Create time</StatLabel>
    <StatNumber fontSize={variables.bodyFontSize}>
      {prettyDate(file.createTime)}
    </StatNumber>
  </Stat>
)

export default CreateTimeInfo
