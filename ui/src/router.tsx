import { createBrowserRouter } from 'react-router-dom'
import AccountInvitationsPage from '@/pages/account/invitations'
import AccountSettingsPage from '@/pages/account/settings'
import AuthenticatedPage from '@/pages/authenticated'
import ConfirmEmailPage from '@/pages/confirm-email'
import FileViewerPage from '@/pages/file-viewer'
import ForgotPasswordPage from '@/pages/forgot-password'
import GroupListPage from '@/pages/group/list'
import GroupMembersPage from '@/pages/group/members'
import GroupSettingsPage from '@/pages/group/settings'
import NewGroupPage from '@/pages/new-group'
import NewOrganizationPage from '@/pages/new-organization'
import NewWorkspacePage from '@/pages/new-workspace'
import OrganizationInvitationsPage from '@/pages/organization/invitations'
import OrganizationLayout from '@/pages/organization/layout'
import OrganizationListPage from '@/pages/organization/list'
import OrganizationMembersPage from '@/pages/organization/members'
import OrganizationSettingsPage from '@/pages/organization/settings'
import ResetPasswordPage from '@/pages/reset-password'
import Root from '@/pages/root'
import SignInPage from '@/pages/sign-in'
import SignOutPage from '@/pages/sign-out'
import SignUpPage from '@/pages/sign-up'
import WorkspaceFilesPage from '@/pages/workspace/files'
import WorkspaceListPage from '@/pages/workspace/list'
import WorkspaceSettingsPage from '@/pages/workspace/settings'
import AccountLayout from './pages/account/layout'
import GroupLayout from './pages/group/layout'
import WorkspaceLayout from './pages/workspace/layout'

const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    children: [
      {
        element: <AuthenticatedPage />,
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
        path: '/reset-password',
        element: <ResetPasswordPage />,
      },
      {
        path: '/confirm-email',
        element: <ConfirmEmailPage />,
      },
    ],
  },
])

export default router
