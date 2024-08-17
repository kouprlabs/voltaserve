// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Config } from './types'

const config: Config = {
  apiURL: '/proxy/api/v2',
  idpURL: '/proxy/idp/v2',
  adminURL: '/proxy/admin/v1',
}

export function getConfig(): Config {
  return config
}
