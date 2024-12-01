// Copyright (c) 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { Badge, Center, Flex, Grid, GridItem, Heading, Spacer, Table, Tbody, Td, Th, Thead, Tr } from '@chakra-ui/react'
import { IconChevronRight, IconFlag, IconGroup, IconPerson, IconWorkspaces, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import semver from 'semver'
import ConsoleAPI, { ComponentVersion } from '@/client/console/console'

const spinnerHeight = '40px'
const uiCurrentVersion = { version: '3.0.0' }
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
const compareFn = (a: ComponentVersion, b: ComponentVersion) => (a.name > b.name ? 1 : 0)

const ConsolePanelOverview = () => {
  const [userCount, setUserCount] = useState<number>()
  const [groupCount, setGroupCount] = useState<number>()
  const [organizationCount, setOrganizationCount] = useState<number>()
  const [workspaceCount, setWorkspaceCount] = useState<number>()
  const [componentsData, setComponentsData] = useState<ComponentVersion[]>([])

  useEffect(() => {
    ConsoleAPI.countObject('user').then((value) => {
      setUserCount(value.count)
    })
    ConsoleAPI.countObject('organization').then((value) => {
      setOrganizationCount(value.count)
    })
    ConsoleAPI.countObject('group').then((value) => {
      setGroupCount(value.count)
    })
    ConsoleAPI.countObject('workspace').then((value) => {
      setWorkspaceCount(value.count)
    })
    internalComponents.map((component) => {
      ConsoleAPI.getComponentsVersions(component).then((value) => {
        if (component.id == 'ui') {
          value.currentVersion = uiCurrentVersion.version
          value.updateAvailable = semver.gt(value.latestVersion, uiCurrentVersion.version)
        }
        setComponentsData((prevState) => {
          return [...prevState.filter((item) => item.name !== value.name), value].toSorted(compareFn)
        })
      })
    })
  }, [])

  return (
    <>
      <Helmet>
        <title>Console</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Console</Heading>
        <Grid gap={4} templateColumns="repeat(4, 1fr)">
          <GridItem>
            <Table>
              <Thead>
                <Tr>
                  <Th>
                    <span className={cx('font-bold')}>Users</span>
                  </Th>
                </Tr>
              </Thead>
            </Table>
            {userCount ? (
              <Heading>
                <div className={cx('flex', 'items-center', 'gap-1', 'p-2')}>
                  <IconPerson className={cx('text-[26px]')} />
                  {userCount}
                </div>
              </Heading>
            ) : (
              <SectionSpinner height={spinnerHeight} />
            )}
          </GridItem>
          <GridItem>
            <Table>
              <Thead>
                <Tr>
                  <Th>
                    <span className={cx('font-bold')}>Organizations</span>
                  </Th>
                </Tr>
              </Thead>
            </Table>
            {organizationCount ? (
              <Heading>
                <div className={cx('flex', 'items-center', 'gap-1', 'p-2')}>
                  <IconFlag className={cx('text-[26px]')} />
                  {organizationCount}
                </div>
              </Heading>
            ) : (
              <SectionSpinner height={spinnerHeight} />
            )}
          </GridItem>
          <GridItem>
            <Table>
              <Thead>
                <Tr>
                  <Th>
                    <span className={cx('font-bold')}>Workspaces</span>
                  </Th>
                </Tr>
              </Thead>
            </Table>
            {workspaceCount ? (
              <Heading>
                <div className={cx('flex', 'items-center', 'gap-1', 'p-2')}>
                  <IconWorkspaces className={cx('text-[26px]')} />
                  {workspaceCount}
                </div>
              </Heading>
            ) : (
              <SectionSpinner height={spinnerHeight} />
            )}
          </GridItem>
          <GridItem>
            <Table>
              <Thead>
                <Tr>
                  <Th>
                    <span className={cx('font-bold')}>Groups</span>
                  </Th>
                </Tr>
              </Thead>
            </Table>
            {groupCount ? (
              <Heading>
                <div className={cx('flex', 'items-center', 'gap-1', 'p-2')}>
                  <IconGroup className={cx('text-[26px]')} />
                  {groupCount}
                </div>
              </Heading>
            ) : (
              <SectionSpinner height={spinnerHeight} />
            )}
          </GridItem>
          <GridItem colSpan={2}>
            <Table>
              <Thead>
                <Tr>
                  <Th>
                    <Flex padding={1}>
                      Components
                      <Spacer />
                      {componentsData.filter((component) => component.updateAvailable).length > 1 ? (
                        <Badge colorScheme="yellow">Updates available</Badge>
                      ) : componentsData.filter((component) => component.updateAvailable).length === 1 ? (
                        <Badge colorScheme="yellow">Update available</Badge>
                      ) : null}
                    </Flex>
                  </Th>
                </Tr>
              </Thead>
              <Tbody>
                {componentsData.map((component) => (
                  <Tr key={component.name}>
                    <Td>
                      <Flex padding={2}>
                        <Link to={component.location}>
                          <Badge>{component.name}</Badge>
                        </Link>
                        <Spacer />
                        {component.updateAvailable ? (
                          <Link to={component.location}>
                            <Center>
                              <Badge colorScheme="red">{component.currentVersion}</Badge>
                              <IconChevronRight />
                              <Badge colorScheme="green">{component.latestVersion}</Badge>
                            </Center>
                          </Link>
                        ) : (
                          <Badge>{component.currentVersion}</Badge>
                        )}
                      </Flex>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </GridItem>
          <GridItem colSpan={2}></GridItem>
        </Grid>
      </div>
    </>
  )
}

export default ConsolePanelOverview
