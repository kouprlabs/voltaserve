// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { IconCheckCircle } from '@koupr/ui'
import cx from 'classnames'
import { FileViewType } from '@/types/file'

type MultiSelectCheckboxProps = {
  isChecked?: boolean
  viewType: FileViewType
}

const MultiSelectCheckbox = ({
  isChecked,
  viewType,
}: MultiSelectCheckboxProps) => {
  return (
    <div
      className={cx('w-[18px]', 'h-[18px]', {
        'relative': viewType === FileViewType.List,
        'absolute': viewType === FileViewType.Grid,
        'top-0.5': viewType === FileViewType.Grid,
        'left-0.5': viewType === FileViewType.Grid,
      })}
    >
      <div
        className={cx(
          'absolute',
          'top-0',
          'left-0',
          'flex',
          'items-center',
          'justify-center',
          'w-[18px]',
          'h-[18px]',
        )}
      >
        <span
          className={cx(
            'z-10',
            {
              'bg-gray-300': !isChecked,
              'dark:bg-gray-400': !isChecked,
              'w-[15px]': !isChecked,
              'h-[15px]': !isChecked,
            },
            {
              'bg-white': isChecked,
              'w-[12px]': isChecked,
              'h-[12px]': isChecked,
            },
            'rounded-full',
          )}
        ></span>
      </div>
      {isChecked ? (
        <IconCheckCircle
          className={cx(
            'absolute',
            'top-0',
            'left-0',
            'z-20',
            'text-[18px]',
            'leading-[18px]',
            {
              'text-blue-500': isChecked,
              'text-gray-500': !isChecked,
              'dark:text-gray-600': !isChecked,
            },
          )}
          filled={isChecked}
        />
      ) : null}
    </div>
  )
}

export default MultiSelectCheckbox
