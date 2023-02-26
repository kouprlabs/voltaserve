import React, { useState } from 'react'
import { Input, Select, Stack } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { FieldAttributes, FieldProps } from 'formik'
import { Unit, getUnit, convertFromByte, normalizeToByte } from './unit'

const StorageInput = ({ id, field, form }: FieldAttributes<FieldProps>) => {
  const [value, setValue] = useState<number | null>(
    field.value ? convertFromByte(field.value, getUnit(field.value)) : null
  )
  const [unit, setUnit] = useState<Unit>(
    field.value ? getUnit(field.value) : 'b'
  )

  return (
    <>
      <input id={id} type="hidden" {...field} />
      <Stack direction="column" spacing={variables.spacing}>
        <Stack direction="row">
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
                  newUnit
                )
                setValue(newValue)
                form.setFieldValue(
                  field.name,
                  normalizeToByte(newValue, newUnit)
                )
              }
            }}
          >
            <option value="b">B</option>
            <option value="mb">MB</option>
            <option value="gb">GB</option>
            <option value="tb">TB</option>
          </Select>
        </Stack>
      </Stack>
    </>
  )
}

export default StorageInput
