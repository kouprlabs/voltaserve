// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { progressAnatomy as parts } from '@chakra-ui/anatomy'
import {
  mode,
  PartsStyleFunction,
  StyleFunctionProps,
  SystemStyleFunction,
  SystemStyleObject,
} from '@chakra-ui/theme-tools'
import variables from '../../variables'

function filledStyle(props: StyleFunctionProps): SystemStyleObject {
  const { colorScheme, hasStripe } = props
  if (hasStripe) {
    return { bg: variables.gradiant }
  } else {
    return { bgColor: mode(`${colorScheme}.500`, `${colorScheme}.200`)(props) }
  }
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const baseStyleFilledTrack: SystemStyleFunction = (props: any) => {
  return {
    ...filledStyle(props),
  }
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const baseStyle: PartsStyleFunction<typeof parts> = (props: any) => ({
  filledTrack: baseStyleFilledTrack(props),
  track: {
    borderRadius: variables.borderRadius,
  },
})

const progress = {
  baseStyle,
}

export default progress
