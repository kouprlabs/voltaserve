// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { MeiliSearch } from 'meilisearch'
import { getConfig } from '@/config/config'

export const USER_SEARCH_INDEX = 'user'

const client = new MeiliSearch({ host: getConfig().search.url })
client.createIndex(USER_SEARCH_INDEX, { primaryKey: 'id' })

export default client
