// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { mode, StyleFunctionProps } from '@chakra-ui/theme-tools'
import variables from '../../variables'

const tab = {
  variants: {
    'solid-rounded': (props: StyleFunctionProps) => ({
      tab: {
        fontSize: variables.bodyFontSize,
        _focus: {
          boxShadow: 'none',
        },
        _selected: {
          bg: mode('black', 'white')(props),
        },
      },
      tabpanel: {
        p: '60px 0 0 0',
      },
    }),
    'line': {
      tab: {
        fontSize: variables.bodyFontSize,
        _focus: {
          boxShadow: 'none',
        },
      },
    },
  },
}

export default tab
