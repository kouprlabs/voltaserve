import {
  Progress,
  Stack,
  Stat,
  StatLabel,
  StatNumber,
  Text,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'
import StorageAPI from '@/client/api/storage'
import prettyBytes from '@/helpers/pretty-bytes'

type StorageUsageInfoProps = {
  file: File
}

const StorageUsageInfo = ({ file }: StorageUsageInfoProps) => {
  const { data, error } = StorageAPI.useGetFileUsage(file.id)
  return (
    <Stat>
      <StatLabel>Storage usage</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        <Stack spacing={variables.spacingXs}>
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
        </Stack>
      </StatNumber>
    </Stat>
  )
}

export default StorageUsageInfo
