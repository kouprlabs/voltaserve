import { createBrowserRouter } from 'react-router-dom'
import AccountInvitationsPage from '@/pages/account/account-invitations-page'
import AccountLayout from '@/pages/account/account-layout'
import AccountSettingsPage from '@/pages/account/account-settings-page'
import ConfirmEmailPage from '@/pages/confirm-email-page'
import FileViewerPage from '@/pages/file-viewer-page'
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
import WorkspaceFilesPage from '@/pages/workspace/workspace-files-page'
import WorkspaceLayout from '@/pages/workspace/workspace-layout'
import WorkspaceListPage from '@/pages/workspace/workspace-list-page'
import WorkspaceSettingsPage from '@/pages/workspace/workspace-settings-page'
import LayoutShell from './components/layout/layout-shell'
import UpdateEmailPage from './pages/update-email-page'

const router = createBrowserRouter([
  {
    path: '/',
    element: <RootPage />,
    children: [
      {
        element: <LayoutShell />,
        children: [
          {
            element: <AccountLayout />,
            children: [
              {
                path: '/account/settings',
                element: <AccountSettingsPage />,
              },
              {
                path: '/account/invitation',
                element: <AccountInvitationsPage />,
              },
            ],
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
        element: <FileViewerPage />,
      },
      {
        path: '/file/:id/mosaic',
        element: <FileViewerPage />,
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
    ],
  },
])

export default router
