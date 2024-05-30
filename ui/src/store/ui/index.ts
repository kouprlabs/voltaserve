import { combineReducers } from 'redux'
import account from './account'
import error from './error'
import files from './files'
import group from './group'
import groupMembers from './group-members'
import groups from './groups'
import incomingInvitations from './incoming-invitations'
import insights from './insights'
import mosaic from './mosaic'
import nav from './nav'
import notifications from './notifications'
import organization from './organization'
import organizations from './organizations'
import outgoingInvitations from './outgoing-invitations'
import snapshots from './snapshots'
import uploadsDrawer from './uploads-drawer'
import watermark from './watermark'
import workspace from './workspace'
import workspaces from './workspaces'

export default combineReducers({
  insights,
  mosaic,
  watermark,
  uploadsDrawer,
  files,
  snapshots,
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
