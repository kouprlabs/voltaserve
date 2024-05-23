import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import { IconDelete, IconSync } from '@/lib'

const InsightsOverviewSettings = () => {
  return (
    <div className={cx('flex', 'flex-row', 'items-stretch', 'gap-1.5')}>
      <Card size="md" variant="outline">
        <CardBody>
          <Text>
            Recollect insights data using the active snapshot, the defined
            language is kept unchanged.
          </Text>
        </CardBody>
        <CardFooter>
          <Button leftIcon={<IconSync />}>Recollect Insights</Button>
        </CardFooter>
      </Card>
      <Card size="md" variant="outline">
        <CardBody>
          <Text>
            Delete insights data from all snapshots, the defined language is
            reset to an empty value.
          </Text>
        </CardBody>
        <CardFooter>
          <Button colorScheme="red" leftIcon={<IconDelete />}>
            Delete Insights
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}

export default InsightsOverviewSettings
