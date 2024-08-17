// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  Badge,
  Center,
  Divider,
  Flex,
  Grid,
  GridItem,
  Heading,
  Spacer,
  Stack,
  StackItem,
  Text,
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import AdminApi, { ComponentVersionList } from '@/client/admin/admin'
import {
  IconChevronRight,
  IconFlag,
  IconGroup,
  IconPerson,
  IconWorkspaces,
} from '@/lib/components/icons'
import SectionSpinner from '@/lib/components/section-spinner'

const spinnerHeight = '40px'

const AdminPanelOverview = () => {
  const [usersAmount, setUsersAmount] = useState<number | undefined>(undefined)
  const [groupsAmount, setGroupsAmount] = useState<number | undefined>(
    undefined,
  )
  const [organizationsAmount, setOrganizationsAmount] = useState<
    number | undefined
  >(undefined)
  const [workspacesAmount, setWorkspacesAmount] = useState<number | undefined>(
    undefined,
  )
  const [componentsData, setComponentsData] = useState<
    ComponentVersionList | undefined
  >(undefined)

  useEffect(() => {
    AdminApi.countObject('user').then((value) => {
      setUsersAmount(value.count)
    })
    AdminApi.countObject('organization').then((value) => {
      setOrganizationsAmount(value.count)
    })
    AdminApi.countObject('group').then((value) => {
      setGroupsAmount(value.count)
    })
    AdminApi.countObject('workspace').then((value) => {
      setWorkspacesAmount(value.count)
    })
    AdminApi.getComponentsVersions().then((value) => {
      setComponentsData(value)
    })
  }, [])

  return (
    <>
      <Helmet>
        <title>Admin Panel</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Cloud Panel</Heading>
        <Grid gap={4} templateColumns="repeat(4, 1fr)">
          <GridItem>
            <Text>
              <span className={cx('font-bold')}>Users</span>
              <Divider />
            </Text>
            {usersAmount ? (
              <Heading>
                <IconPerson className={cx('text-[26px]')} />
                {usersAmount}
              </Heading>
            ) : (
              <SectionSpinner height={spinnerHeight} />
            )}
          </GridItem>
          <GridItem>
            <Text>
              <span className={cx('font-bold')}>Organizations</span>
              <Divider />
            </Text>
            {organizationsAmount ? (
              <Heading>
                <IconFlag className={cx('text-[26px]')} />
                {organizationsAmount}
              </Heading>
            ) : (
              <SectionSpinner height={spinnerHeight} />
            )}
          </GridItem>
          <GridItem>
            <Text>
              <span className={cx('font-bold')}>Workspaces</span>
              <Divider />
            </Text>
            {workspacesAmount ? (
              <Heading>
                <IconWorkspaces className={cx('text-[26px]')} />
                {workspacesAmount}
              </Heading>
            ) : (
              <SectionSpinner height={spinnerHeight} />
            )}
          </GridItem>
          <GridItem>
            <Text>
              <span className={cx('font-bold')}>Groups</span>
              <Divider />
            </Text>
            {groupsAmount ? (
              <Heading>
                <IconGroup className={cx('text-[26px]')} />
                {groupsAmount}
              </Heading>
            ) : (
              <SectionSpinner height={spinnerHeight} />
            )}
          </GridItem>
          <GridItem colSpan={2}>
            <Flex padding={1}>
              <Text>
                <span className={cx('font-bold')}>Statistics</span>
              </Text>
              <Spacer />
            </Flex>
            <Divider />
          </GridItem>
          {componentsData ? (
            <GridItem colSpan={2}>
              <Flex padding={1}>
                <Text>
                  <span className={cx('font-bold')}>Components</span>
                </Text>
                <Spacer />
                {componentsData.data.filter(
                  (component) => component.updateAvailable,
                ).length > 1 ? (
                  <Badge colorScheme="yellow">Updates available</Badge>
                ) : componentsData.data.filter(
                    (component) => component.updateAvailable,
                  ).length === 1 ? (
                  <Badge colorScheme="yellow">Update available</Badge>
                ) : null}
              </Flex>
              <Divider />
              <Stack>
                {componentsData.data.map((component) => (
                  <>
                    <StackItem>
                      <Flex padding={2}>
                        <Text>{component.name}</Text>
                        <Spacer />
                        {component.updateAvailable ? (
                          <Link to={component.location}>
                            <Center>
                              <Badge colorScheme="red">
                                {component.currentVersion}
                              </Badge>
                              <IconChevronRight />
                              <Badge colorScheme="green">
                                {component.latestVersion}
                              </Badge>
                            </Center>
                          </Link>
                        ) : (
                          <Badge>{component.currentVersion}</Badge>
                        )}
                      </Flex>
                    </StackItem>
                    <Divider />
                  </>
                ))}
              </Stack>
            </GridItem>
          ) : (
            <SectionSpinner />
          )}
        </Grid>
      </div>
    </>
  )
}

export default AdminPanelOverview
