// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { combineReducers } from 'redux'
import account from './account'
import error from './error'
import files from './files'
import group from './group'
import groupMembers from './group-members'
import groups from './groups'
import incomingInvitations from './incoming-invitations'
import indexes from './indexes'
import insights from './insights'
import mosaic from './mosaic'
import nav from './nav'
import organization from './organization'
import organizations from './organizations'
import outgoingInvitations from './outgoing-invitations'
import searchFilter from './search-filter'
import snapshots from './snapshots'
import notifications from './tasks'
import tasks from './tasks'
import uploads from './uploads'
import workspace from './workspace'
import workspaces from './workspaces'

export default combineReducers({
  insights,
  mosaic,
  searchFilter,
  uploads,
  tasks,
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
  indexes,
  account,
  notifications,
  workspace,
  workspaces,
})
