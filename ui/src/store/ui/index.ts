import { combineReducers } from 'redux'
import account from './account'
import ai from './ai'
import error from './error'
import files from './files'
import group from './group'
import groupMembers from './group-members'
import groups from './groups'
import incomingInvitations from './incoming-invitations'
import nav from './nav'
import notifications from './notifications'
import organization from './organization'
import organizations from './organizations'
import outgoingInvitations from './outgoing-invitations'
import uploadsDrawer from './uploads-drawer'
import workspace from './workspace'
import workspaces from './workspaces'

export default combineReducers({
  ai,
  uploadsDrawer,
  files,
  nav,
  error,
  organization,
  organizations,
  group,
  groupMembers,
  groups,
  outgoingInvitations,
  incomingInvitations,
  account,
  notifications,
  workspace,
  workspaces,
})
