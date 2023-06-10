import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom'
import { Heading, Stack, Tab, TabList, Tabs } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import GroupAPI from '@/api/group'
import { swrConfig } from '@/api/options'
import { geEditorPermission } from '@/api/permission'

const GroupLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const params = useParams()
  const groupId = params.id as string
  const { data: group } = GroupAPI.useGetById(groupId, swrConfig())
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
    <Stack direction="column" spacing={variables.spacing2Xl}>
      <Heading size="lg">{group.name}</Heading>
      <Tabs variant="solid-rounded" colorScheme="gray" index={tabIndex}>
        <TabList>
          <Tab onClick={() => navigate(`/group/${groupId}/member`)}>
            Members
          </Tab>
          {geEditorPermission(group.permission) && (
            <Tab onClick={() => navigate(`/group/${groupId}/settings`)}>
              Settings
            </Tab>
          )}
        </TabList>
      </Tabs>
      <Outlet />
    </Stack>
  )
}

export default GroupLayout
