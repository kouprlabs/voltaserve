import { combineReducers } from 'redux'
import error from './error'
import files from './files'
import groupMembers from './group-members'
import nav from './nav'
import organizations from './organizations'
import outgoingInvitations from './outgoing-invitations'
import uploadsDrawer from './uploads-drawer'

export default combineReducers({
  uploadsDrawer,
  files,
  nav,
  error,
  organizations,
  groupMembers,
  outgoingInvitations,
})
