import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom'
import { Heading, Stack, Tab, TabList, Tabs } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'

const WorkspaceLayout = () => {
  const location = useLocation()
  const params = useParams()
  const navigate = useNavigate()
  const workspaceId = params.id as string
  const { data: workspace } = WorkspaceAPI.useGetById(workspaceId, swrConfig())
  const [tabIndex, setTabIndex] = useState(0)

  useEffect(() => {
    const segments = location.pathname.split('/')
    const segment = segments[segments.length - 1]
    if (segment === 'settings') {
      setTabIndex(1)
    } else {
      setTabIndex(0)
    }
  }, [location])

  if (!workspace) {
    return null
  }

  return (
    <Stack direction="column" spacing={variables.spacing2Xl} height="100%">
      <Heading size="lg">{workspace.name}</Heading>
      <Tabs variant="solid-rounded" colorScheme="gray" index={tabIndex}>
        <TabList>
          <Tab
            onClick={() =>
              navigate(`/workspace/${workspaceId}/file/${workspace.rootId}`)
            }
          >
            Files
          </Tab>
          <Tab onClick={() => navigate(`/workspace/${workspaceId}/settings`)}>
            Settings
          </Tab>
        </TabList>
      </Tabs>
      <Outlet />
    </Stack>
  )
}

export default WorkspaceLayout
