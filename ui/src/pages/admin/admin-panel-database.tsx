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
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { Heading, Tab, TabList, Tabs } from '@chakra-ui/react'
import cx from 'classnames'
import AdminApi from '@/client/admin/admin'

const AdminPanelDatabase = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const [tabIndex, setTabIndex] = useState(0)
  const [indexesAvailable, setIndexesAvailable] = useState(true)

  useEffect(() => {
    const segments = location.pathname.split('/')
    const segment = segments[segments.length - 1]
    if (segment === 'overview') {
      setTabIndex(0)
    } else if (segment === 'indexes') {
      setTabIndex(1)
    }
  }, [location])

  useEffect(() => {
    AdminApi.checkIndexesAvailability().then((value) => {
      setIndexesAvailable(value)
    })
  }, [])

  return (
    <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
      <Heading className={cx('text-heading', 'shrink-0')} noOfLines={1}>
        Database management
      </Heading>
      <Tabs variant="solid-rounded" colorScheme="gray" index={tabIndex}>
        <TabList>
          <Tab onClick={() => navigate(`/admin/database/overview`)}>
            Overview
          </Tab>
          {indexesAvailable ? (
            <Tab onClick={() => navigate(`/admin/database/indexes`)}>
              Indexes
            </Tab>
          ) : null}
        </TabList>
      </Tabs>
      <Outlet />
    </div>
  )
}

export default AdminPanelDatabase
