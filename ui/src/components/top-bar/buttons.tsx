import { Link } from 'react-router-dom'
import { Button } from '@chakra-ui/react'
import { IconAdd } from '@/components/common/icon'

export const CreateGroupButton = () => (
  <Button
    as={Link}
    to="/new/group"
    leftIcon={<IconAdd />}
    variant="solid"
    colorScheme="blue"
  >
    New group
  </Button>
)

export const CreateOrganizationButton = () => (
  <Button
    as={Link}
    to="/new/organization"
    leftIcon={<IconAdd />}
    variant="solid"
    colorScheme="blue"
  >
    New organization
  </Button>
)

export const CreateWorkspaceButton = () => (
  <Button
    as={Link}
    to="/new/workspace"
    leftIcon={<IconAdd />}
    variant="solid"
    colorScheme="blue"
  >
    New workspace
  </Button>
)
