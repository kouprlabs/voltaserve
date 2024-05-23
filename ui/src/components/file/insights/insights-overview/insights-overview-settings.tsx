import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import { IconDelete, IconSync } from '@/lib'

const InsightsOverviewSettings = () => {
  return (
    <div className={cx('flex', 'flex-row', 'items-stretch', 'gap-1.5')}>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>
            Creates new insights for the active snapshot, uses the existing
            language.
          </Text>
        </CardBody>
        <CardFooter>
          <Button leftIcon={<IconSync />}>Update</Button>
        </CardFooter>
      </Card>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>
            Deletes insights from the active snapshot, can be recreated later.
          </Text>
        </CardBody>
        <CardFooter>
          <Button colorScheme="red" leftIcon={<IconDelete />}>
            Delete
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}

export default InsightsOverviewSettings
