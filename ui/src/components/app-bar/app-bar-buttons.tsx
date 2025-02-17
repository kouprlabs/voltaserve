// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Link } from 'react-router-dom'
import { Button } from '@chakra-ui/react'
import { IconAdd } from '@koupr/ui'

export const CreateGroupButton = () => (
  <Button
    as={Link}
    to="/new/group"
    leftIcon={<IconAdd />}
    variant="solid"
    colorScheme="blue"
  >
    New Group
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
    New Organization
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
    New Workspace
  </Button>
)
