// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import variables from '../../variables'

const checkbox = {
  baseStyle: {
    control: {
      borderRadius: '50%',
    },
  },
  sizes: {
    md: {
      control: { w: '20px', h: '20px' },
      label: {
        fontSize: variables.bodyFontSize,
      },
    },
  },
}

export default checkbox
