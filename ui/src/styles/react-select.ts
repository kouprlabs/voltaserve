// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { ColorMode, SystemStyleObject } from '@chakra-ui/system'

export type ReactSelectStylesOptions = {
  colorMode?: ColorMode
}

export function reactSelectStyles(options?: ReactSelectStylesOptions) {
  let bg = 'transparent'
  if (options?.colorMode === 'light') {
    bg = 'white'
  } else if (options?.colorMode === 'dark') {
    bg = 'gray.800'
  }
  return {
    control: (provided: SystemStyleObject) => ({
      ...provided,
      bg,
    }),
    dropdownIndicator: (provided: SystemStyleObject) => ({
      ...provided,
      bg,
      cursor: 'inherit',
      position: 'absolute',
      right: '0px',
    }),
    menuList: (provided: SystemStyleObject) => ({
      ...provided,
      borderRadius: '15px',
    }),
    indicatorSeparator: (provided: SystemStyleObject) => ({
      ...provided,
      display: 'none',
    }),
    placeholder: (provided: SystemStyleObject) => ({
      ...provided,
      textAlign: 'center',
    }),
    singleValue: (provided: SystemStyleObject) => ({
      ...provided,
      textAlign: 'center',
    }),
  }
}

export default reactSelectStyles
