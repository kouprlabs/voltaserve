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
      <div className={cx('flex', 'flex-col', 'gap-1.5')}>
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
          <Select
            defaultValue={unit}
            className={cx('shrink-0', 'w-auto')}
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
