import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom'
import { Heading, Tab, TabList, Tabs } from '@chakra-ui/react'
import cx from 'classnames'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/workspace'

const WorkspaceLayout = () => {
  const location = useLocation()
  const { id } = useParams()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const { data: workspace, mutate } = WorkspaceAPI.useGetById(id, swrConfig())
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

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  if (!workspace) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-col', 'gap-2', 'h-full')}>
      <Heading className={cx('text-heading')}>{workspace.name}</Heading>
      <Tabs variant="solid-rounded" colorScheme="gray" index={tabIndex}>
        <TabList>
          <Tab
            onClick={() =>
              navigate(`/workspace/${id}/file/${workspace.rootId}`)
            }
          >
            Files
          </Tab>
          <Tab onClick={() => navigate(`/workspace/${id}/settings`)}>
            Settings
          </Tab>
        </TabList>
      </Tabs>
      <Outlet />
    </div>
  )
}

export default WorkspaceLayout
