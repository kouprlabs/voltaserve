import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { File } from '@/api/file'
import variables from '@/theme/variables'
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
