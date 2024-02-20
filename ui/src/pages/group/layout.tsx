import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom'
import { Heading, Tab, TabList, Tabs } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import classNames from 'classnames'
import GroupAPI from '@/client/api/group'
import { swrConfig } from '@/client/options'

const GroupLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { id } = useParams()
  const { data: group } = GroupAPI.useGetById(id, swrConfig())
  const [tabIndex, setTabIndex] = useState(0)

  useEffect(() => {
    const segments = location.pathname.split('/')
    const segment = segments[segments.length - 1]
    if (segment === 'member') {
      setTabIndex(0)
    } else if (segment === 'settings') {
      setTabIndex(1)
    }
  }, [location])

  if (!group) {
    return null
  }

  return (
    <div className={classNames('flex', 'flex-col', 'gap-3.5')}>
      <Heading fontSize={variables.headingFontSize}>{group.name}</Heading>
      <Tabs variant="solid-rounded" colorScheme="gray" index={tabIndex}>
        <TabList>
          <Tab onClick={() => navigate(`/group/${id}/member`)}>Members</Tab>
          <Tab onClick={() => navigate(`/group/${id}/settings`)}>Settings</Tab>
        </TabList>
      </Tabs>
      <Outlet />
    </div>
  )
}

export default GroupLayout
