import { Progress, Stat, StatLabel, StatNumber, Text } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import classNames from 'classnames'
import { File } from '@/client/api/file'
import StorageAPI from '@/client/api/storage'
import prettyBytes from '@/helpers/pretty-bytes'

export type FileInfoStorageUsageProps = {
  file: File
}

const FileInfoStorageUsage = ({ file }: FileInfoStorageUsageProps) => {
  const { data, error } = StorageAPI.useGetFileUsage(file.id)
  return (
    <Stat>
      <StatLabel>Storage usage</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        <div className={classNames('flex', 'flex-col', 'gap-0.5')}>
          {error && <Text color="red">Failed to load.</Text>}
          {data && !error && (
            <>
              <Text>
                {prettyBytes(data.bytes)} of {prettyBytes(data.maxBytes)} used
              </Text>
              <Progress size="sm" value={data.percentage} hasStripe />
            </>
          )}
          {!data && !error && (
            <>
              <Text>Calculating…</Text>
              <Progress size="sm" value={0} hasStripe />
            </>
          )}
        </div>
      </StatNumber>
    </Stat>
  )
}

export default FileInfoStorageUsage
