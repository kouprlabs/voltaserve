import { Badge } from '@chakra-ui/react'
import { InvitationStatus } from '@/client/api/invitation'

type OrganizationStatusProps = {
  value: InvitationStatus
}

const OrganizationStatus = ({ value }: OrganizationStatusProps) => {
  let colorScheme
  if (value === 'accepted') {
    colorScheme = 'green'
  } else if (value === 'declined') {
    colorScheme = 'red'
  }
  return <Badge colorScheme={colorScheme}>{value}</Badge>
}

export default OrganizationStatus
