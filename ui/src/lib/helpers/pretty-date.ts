// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export default function prettyDate(date: any): string {
  return new Date(date).toLocaleString()
}
