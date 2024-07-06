// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { useState } from 'react'
import { Input, Select } from '@chakra-ui/react'
import { FieldAttributes, FieldProps } from 'formik'
import cx from 'classnames'
import { Unit, convertFromByte, getUnit, normalizeToByte } from './convert'

const StorageInput = ({ id, field, form }: FieldAttributes<FieldProps>) => {
  const [value, setValue] = useState<number | null>(
    field.value ? convertFromByte(field.value, getUnit(field.value)) : null,
  )
  const [unit, setUnit] = useState<Unit>(
    field.value ? getUnit(field.value) : 'b',
  )

  return (
    <>
      <input id={id} type="hidden" {...field} />
      <div className={cx('flex', 'flex-row', 'gap-0.5')}>
        <Input
          type="number"
          disabled={form.isSubmitting}
          value={value || ''}
          onChange={(event) => {
            if (event.target.value) {
              const newValue = parseInt(event.target.value)
              setValue(newValue)
              form.setFieldValue(field.name, normalizeToByte(newValue, unit))
            } else {
              setValue(null)
              form.setFieldValue(field.name, '')
            }
          }}
        />
        <div className={cx('min-w-[80px]')}>
          <Select
            defaultValue={unit}
            disabled={form.isSubmitting}
            onChange={(event) => {
              const newUnit = event.target.value as Unit
              setUnit(newUnit)
              if (value) {
                const newValue = convertFromByte(
                  normalizeToByte(value, unit),
                  newUnit,
                )
                setValue(newValue)
                form.setFieldValue(
                  field.name,
                  normalizeToByte(newValue, newUnit),
                )
              }
            }}
          >
            <option value="b">B</option>
            <option value="mb">MB</option>
            <option value="gb">GB</option>
            <option value="tb">TB</option>
          </Select>
        </div>
      </div>
    </>
  )
}

export default StorageInput
