import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom'
import { Heading, Stack, Tab, TabList, Tabs } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { swrConfig } from '@/api/options'
import OrganizationAPI from '@/api/organization'
import { geEditorPermission, geOwnerPermission } from '@/api/permission'

const OrganizationLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const params = useParams()
  const orgId = params.id as string
  const { data: org } = OrganizationAPI.useGetById(orgId, swrConfig())
  const [tabIndex, setTabIndex] = useState(0)

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

  if (!org) {
    return null
  }

  return (
    <Stack direction="column" spacing={variables.spacing2Xl}>
      <Heading size="lg">{org.name}</Heading>
      <Tabs variant="solid-rounded" colorScheme="gray" index={tabIndex}>
        <TabList>
          <Tab onClick={() => navigate(`/organization/${orgId}/member`)}>
            Members
          </Tab>
          {geOwnerPermission(org.permission) && (
            <Tab onClick={() => navigate(`/organization/${orgId}/invitation`)}>
              Invitations
            </Tab>
          )}
          {geEditorPermission(org.permission) && (
            <Tab onClick={() => navigate(`/organization/${orgId}/settings`)}>
              Settings
            </Tab>
          )}
        </TabList>
      </Tabs>
      <Outlet />
    </Stack>
  )
}

export default OrganizationLayout
