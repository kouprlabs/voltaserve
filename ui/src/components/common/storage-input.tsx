import { useState } from 'react'
import { Input, Select } from '@chakra-ui/react'
import { FieldAttributes, FieldProps } from 'formik'
import classNames from 'classnames'
import {
  byteToGigabyte,
  byteToMegabyte,
  byteToTerabyte,
  gigabyteToByte,
  megabyteToByte,
  terabyteToByte,
} from '@/helpers/convert-storage'

type Unit = 'b' | 'mb' | 'gb' | 'tb'

function getUnit(value: number): Unit {
  if (value >= 1e12) {
    return 'tb'
  }
  if (value >= 1e9) {
    return 'gb'
  }
  if (value >= 1e6) {
    return 'mb'
  }
  return 'b'
}

function convertFromByte(value: number, unit: Unit): number {
  if (unit === 'b') {
    return value
  }
  if (unit === 'mb') {
    return byteToMegabyte(value)
  }
  if (unit === 'gb') {
    return byteToGigabyte(value)
  }
  if (unit === 'tb') {
    return byteToTerabyte(value)
  }
  throw new Error(`Invalid unit: ${unit}`)
}

function normalizeToByte(value: number, unit: Unit) {
  if (unit === 'b') {
    return value
  }
  if (unit === 'mb') {
    return megabyteToByte(value)
  }
  if (unit === 'gb') {
    return gigabyteToByte(value)
  }
  if (unit === 'tb') {
    return terabyteToByte(value)
  }
  throw new Error(`Invalid unit: ${unit}`)
}

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
      <div className={classNames('flex', 'flex-col', 'gap-1.5')}>
        <div className={classNames('flex', 'flex-row', 'gap-0.5')}>
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
            flexShrink={0}
            w="auto"
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
