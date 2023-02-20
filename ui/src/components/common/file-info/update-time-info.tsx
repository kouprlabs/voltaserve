import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { File } from '@/api/file'
import variables from '@/theme/variables'
import prettyDate from '@/helpers/pretty-date'

type UpdateTimeInfoProps = {
  file: File
}

const UpdateTimeInfo = ({ file }: UpdateTimeInfoProps) => {
  if (
    !file.updateTime ||
    (file.updateTime &&
      file.createTime.includes(file.updateTime.replaceAll('Z', '')))
  ) {
    return null
  }
  return (
    <Stat>
      <StatLabel>Update time</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        {prettyDate(file.updateTime)}
      </StatNumber>
    </Stat>
  )
}

export default UpdateTimeInfo
