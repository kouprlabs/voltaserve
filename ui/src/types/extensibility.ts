// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { FormSection } from '@koupr/ui'

export type Page = {
  path: string
  component: React.ReactNode
  menu?: Tab
  tab?: Menu
}

export type Tab = {
  label: string
}

export type Menu = {
  label: string
}

export type Extensions = {
  account?: AccountExtensions
}

export type AccountExtensions = {
  pages?: Page[]
  settings?: AccountSettingsExtensions
}

export type AccountSettingsExtensions = {
  sections?: FormSection[]
}
