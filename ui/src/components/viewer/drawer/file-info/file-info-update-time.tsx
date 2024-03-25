import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'
import prettyDate from '@/helpers/pretty-date'

type FileInfoUpdateTimeProps = {
  file: File
}

const FileInfoUpdateTime = ({ file }: FileInfoUpdateTimeProps) => {
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

export default FileInfoUpdateTime
