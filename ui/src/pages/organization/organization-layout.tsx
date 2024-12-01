// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom'
import { Heading, Tab, TabList, Tabs } from '@chakra-ui/react'
import { SectionError, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI from '@/client/api/organization'
import { geOwnerPermission } from '@/client/api/permission'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/organization'

const OrganizationLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const { id } = useParams()
  const {
    data: org,
    error: orgError,
    isLoading: orgIsLoading,
    mutate,
  } = OrganizationAPI.useGet(id, swrConfig())
  const [tabIndex, setTabIndex] = useState(0)
  const orgIsReady = org && !orgError

  useEffect(() => {
    const segments = location.pathname.split('/')
    const segment = segments[segments.length - 1]
    if (segment === 'member') {
      setTabIndex(0)
    } else if (segment === 'invitation') {
      setTabIndex(1)
    } else if (segment === 'settings') {
      setTabIndex(2)
    }
  }, [location])

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  return (
    <>
      {orgIsLoading ? <SectionSpinner /> : null}
      {orgError ? <SectionError text={errorToString(orgError)} /> : null}
      {orgIsReady ? (
        <>
          <Helmet>
            <title>{org.name}</title>
          </Helmet>
          <div className={cx('flex', 'flex-col', 'gap-3.5')}>
            <Heading className={cx('text-heading', 'shrink-0')} noOfLines={1}>
              {org.name}
            </Heading>
            <Tabs variant="solid-rounded" colorScheme="gray" index={tabIndex}>
              <TabList>
                <Tab onClick={() => navigate(`/organization/${id}/member`)}>
                  Members
                </Tab>
                <Tab
                  onClick={() => navigate(`/organization/${id}/invitation`)}
                  display={geOwnerPermission(org.permission) ? 'auto' : 'none'}
                >
                  Invitations
                </Tab>
                <Tab onClick={() => navigate(`/organization/${id}/settings`)}>
                  Settings
                </Tab>
              </TabList>
            </Tabs>
            <Outlet />
          </div>
        </>
      ) : null}
    </>
  )
}

export default OrganizationLayout
