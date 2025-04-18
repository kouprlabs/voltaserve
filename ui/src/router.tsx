// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { createBrowserRouter, RouteObject } from 'react-router-dom'
import LayoutConsole from '@/components/layout/layout-console'
import LayoutShell from '@/components/layout/layout-shell'
import AccountInvitationsPage from '@/pages/account/account-invitations-page'
import AccountLayout from '@/pages/account/account-layout'
import AccountSettingsPage from '@/pages/account/account-settings-page'
import ConfirmEmailPage from '@/pages/confirm-email-page'
import ConsolePanelGroups from '@/pages/console/console-panel-groups'
import ConsolePanelOrganization from '@/pages/console/console-panel-organization'
import ConsolePanelOrganizations from '@/pages/console/console-panel-organizations'
import ConsolePanelOverview from '@/pages/console/console-panel-overview'
import ConsolePanelUser from '@/pages/console/console-panel-user'
import ConsolePanelUsers from '@/pages/console/console-panel-users'
import ConsolePanelWorkspaces from '@/pages/console/console-panel-workspaces'
import ForgotPasswordPage from '@/pages/forgot-password-page'
import GroupLayout from '@/pages/group/group-layout'
import GroupListPage from '@/pages/group/group-list-page'
import GroupMembersPage from '@/pages/group/group-members-page'
import GroupSettingsPage from '@/pages/group/group-settings-page'
import NewGroupPage from '@/pages/new-group-page'
import NewOrganizationPage from '@/pages/new-organization-page'
import NewWorkspacePage from '@/pages/new-workspace-page'
import OrganizationInvitationsPage from '@/pages/organization/organization-invitations-page'
import OrganizationLayout from '@/pages/organization/organization-layout'
import OrganizationListPage from '@/pages/organization/organization-list-page'
import OrganizationMembersPage from '@/pages/organization/organization-members-page'
import OrganizationSettingsPage from '@/pages/organization/organization-settings-page'
import ResetPasswordPage from '@/pages/reset-password-page'
import RootPage from '@/pages/root-page'
import SignInPage from '@/pages/sign-in-page'
import SignOutPage from '@/pages/sign-out-page'
import SignUpPage from '@/pages/sign-up-page'
import UpdateEmailPage from '@/pages/update-email-page'
import ViewerPage from '@/pages/viewer-page'
import WorkspaceFilesPage from '@/pages/workspace/workspace-files-page'
import WorkspaceLayout from '@/pages/workspace/workspace-layout'
import WorkspaceListPage from '@/pages/workspace/workspace-list-page'
import WorkspaceSettingsPage from '@/pages/workspace/workspace-settings-page'
import { Extensions } from '@/types/extensibility'

export type CreateRouterOptions = {
  extensions?: Extensions
}

export function createRouter({ extensions }: CreateRouterOptions) {
  const accountChildren: RouteObject[] = [
    {
      path: '/account/settings',
      element: <AccountSettingsPage />,
    },
    {
      path: '/account/invitation',
      element: <AccountInvitationsPage />,
    },
  ]
  if (extensions?.account?.pages) {
    extensions?.account?.pages.forEach((tab) => {
      accountChildren.push({
        path: tab.path,
        element: tab.component,
      })
    })
  }
  return createBrowserRouter([
    {
      path: '/',
      element: <RootPage />,
      children: [
        {
          element: <LayoutShell extensions={extensions} />,
          children: [
            {
              element: <AccountLayout extensions={extensions?.account} />,
              children: accountChildren,
            },
            {
              path: '/workspace',
              element: <WorkspaceListPage />,
            },
            {
              element: <WorkspaceLayout />,
              children: [
                {
                  path: '/workspace/:id/file/:fileId',
                  element: <WorkspaceFilesPage />,
                },
                {
                  path: '/workspace/:id/settings',
                  element: <WorkspaceSettingsPage />,
                },
              ],
            },
            {
              path: '/organization',
              element: <OrganizationListPage />,
            },
            {
              element: <OrganizationLayout />,
              children: [
                {
                  path: '/organization/:id/invitation',
                  element: <OrganizationInvitationsPage />,
                },
                {
                  path: '/organization/:id/member',
                  element: <OrganizationMembersPage />,
                },
                {
                  path: '/organization/:id/settings',
                  element: <OrganizationSettingsPage />,
                },
              ],
            },
            {
              path: '/group',
              element: <GroupListPage />,
            },
            {
              element: <GroupLayout />,
              children: [
                {
                  path: '/group/:id/member',
                  element: <GroupMembersPage />,
                },
                {
                  path: '/group/:id/settings',
                  element: <GroupSettingsPage />,
                },
              ],
            },
            {
              path: '/new/workspace',
              element: <NewWorkspacePage />,
            },
            {
              path: '/new/group',
              element: <NewGroupPage />,
            },
            {
              path: '/new/organization',
              element: <NewOrganizationPage />,
            },
          ],
        },
        {
          path: '/file/:id',
          element: <ViewerPage />,
        },
        {
          path: '/file/:id/mosaic',
          element: <ViewerPage />,
        },
        {
          path: '/sign-up',
          element: <SignUpPage />,
        },
        {
          path: '/sign-out',
          element: <SignOutPage />,
        },
        {
          path: '/sign-in',
          element: <SignInPage />,
        },
        {
          path: '/forgot-password',
          element: <ForgotPasswordPage />,
        },
        {
          path: '/reset-password/:token',
          element: <ResetPasswordPage />,
        },
        {
          path: '/confirm-email/:token',
          element: <ConfirmEmailPage />,
        },
        {
          path: '/update-email/:token',
          element: <UpdateEmailPage />,
        },
        {
          element: <LayoutConsole />,
          children: [
            {
              path: '/console/dashboard',
              element: <ConsolePanelOverview />,
            },
            {
              path: '/console/users',
              element: <ConsolePanelUsers />,
            },
            {
              path: '/console/users/:id',
              element: <ConsolePanelUser />,
            },
            {
              path: '/console/workspaces',
              element: <ConsolePanelWorkspaces />,
            },
            {
              path: '/console/organizations',
              element: <ConsolePanelOrganizations />,
            },
            {
              path: '/console/organizations/:id',
              element: <ConsolePanelOrganization />,
            },
            {
              path: '/console/groups',
              element: <ConsolePanelGroups />,
            },
          ],
        },
      ],
    },
  ])
}
