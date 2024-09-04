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
import semver from 'semver'
import ConsoleApi, { ComponentVersion } from '@/client/console/console'
import {
  IconChevronRight,
  IconFlag,
  IconGroup,
  IconPerson,
  IconWorkspaces,
} from '@/lib/components/icons'
import SectionSpinner from '@/lib/components/section-spinner'

const spinnerHeight = '40px'
const uiCurrentVersion = { version: '2.1.0' }
const internalComponents = [
  { id: 'ui' },
  { id: 'api' },
  { id: 'language' },
  { id: 'webdav' },
  { id: 'idp' },
  { id: 'mosaic' },
  { id: 'console' },
  { id: 'conversion' },
]
const compareFn = (a: ComponentVersion, b: ComponentVersion) =>
  a.name > b.name ? 1 : 0

const ConsolePanelOverview = () => {
  const [usersAmount, setUsersAmount] = useState<number>()
  const [groupsAmount, setGroupsAmount] = useState<number>()
  const [organizationsAmount, setOrganizationsAmount] = useState<number>()
  const [workspacesAmount, setWorkspacesAmount] = useState<number>()
  const [componentsData, setComponentsData] = useState<ComponentVersion[]>([])

  useEffect(() => {
    ConsoleApi.countObject('user').then((value) => {
      setUsersAmount(value.count)
    })
    ConsoleApi.countObject('organization').then((value) => {
      setOrganizationsAmount(value.count)
    })
    ConsoleApi.countObject('group').then((value) => {
      setGroupsAmount(value.count)
    })
    ConsoleApi.countObject('workspace').then((value) => {
      setWorkspacesAmount(value.count)
    })
    internalComponents.map((component) => {
      ConsoleApi.getComponentsVersions(component).then((value) => {
        if (component.id == 'ui') {
          value.currentVersion = uiCurrentVersion.version
          value.updateAvailable = semver.gt(
            value.latestVersion,
            uiCurrentVersion.version,
          )
        }
        setComponentsData((prevState) => {
          return [
            ...prevState.filter((item) => item.name !== value.name),
            value,
          ].toSorted(compareFn)
        })
      })
    })
  }, [])

  return (
    <>
      <Helmet>
        <title>Console Panel</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Cloud Console</Heading>
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
                {componentsData.filter((component) => component.updateAvailable)
                  .length > 1 ? (
                  <Badge colorScheme="yellow">Updates available</Badge>
                ) : componentsData.filter(
                    (component) => component.updateAvailable,
                  ).length === 1 ? (
                  <Badge colorScheme="yellow">Update available</Badge>
                ) : null}
              </Flex>
              <Divider />
              <Stack>
                {componentsData.map((component) => (
                  <StackItem key={component.name}>
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
                    <Divider />
                  </StackItem>
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

export default ConsolePanelOverview
