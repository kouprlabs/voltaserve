// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom'
import { Heading, Tab, TabList, Tabs } from '@chakra-ui/react'
import { SectionError, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/workspace'

const WorkspaceLayout = () => {
  const location = useLocation()
  const { id } = useParams()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const {
    data: workspace,
    error: workspaceError,
    isLoading: isWorkspaceLoading,
    mutate,
  } = WorkspaceAPI.useGet(id, swrConfig())
  const [tabIndex, setTabIndex] = useState(0)
  const isWorkspaceError = !workspace && workspaceError
  const isWorkspaceReady = workspace && !workspaceError

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
    <>
      {isWorkspaceLoading ? <SectionSpinner /> : null}
      {isWorkspaceError ? (
        <SectionError text="Failed to load workspace." />
      ) : null}
      {isWorkspaceReady ? (
        <>
          <Helmet>
            <title>{workspace.name}</title>
          </Helmet>
          <div className={cx('flex', 'flex-col', 'gap-2', 'h-full')}>
            <Heading className={cx('text-heading', 'shrink-0')} noOfLines={1}>
              {workspace.name}
            </Heading>
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
        </>
      ) : null}
    </>
  )
}

export default WorkspaceLayout
