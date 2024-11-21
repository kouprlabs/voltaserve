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
import GroupAPI from '@/client/api/group'
import { swrConfig } from '@/client/options'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/group'

const GroupLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const { id } = useParams()
  const {
    data: group,
    error: groupError,
    isLoading: isGroupLoading,
    mutate,
  } = GroupAPI.useGet(id, swrConfig())
  const [tabIndex, setTabIndex] = useState(0)
  const isGroupError = !group && groupError
  const isGroupReady = group && !groupError

  useEffect(() => {
    const segments = location.pathname.split('/')
    const segment = segments[segments.length - 1]
    if (segment === 'member') {
      setTabIndex(0)
    } else if (segment === 'settings') {
      setTabIndex(1)
    }
  }, [location])

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  return (
    <>
      {isGroupLoading ? <SectionSpinner /> : null}
      {isGroupError ? <SectionError text="Failed to load group." /> : null}
      {isGroupReady ? (
        <>
          <Helmet>
            <title>{group.name}</title>
          </Helmet>
          <div className={cx('flex', 'flex-col', 'gap-3.5')}>
            <Heading className={cx('text-heading', 'shrink-0')} noOfLines={1}>
              {group.name}
            </Heading>
            <Tabs variant="solid-rounded" colorScheme="gray" index={tabIndex}>
              <TabList>
                <Tab onClick={() => navigate(`/group/${id}/member`)}>
                  Members
                </Tab>
                <Tab onClick={() => navigate(`/group/${id}/settings`)}>
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

export default GroupLayout
