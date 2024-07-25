import cx from 'classnames'
import { PasswordRequirements } from '@/client/idp/account'
import { IconCheck } from '@/lib/components/icons'

export type PasswordHintsProps = {
  value: string
  requirements: PasswordRequirements
}

const PasswordHints = ({ value, requirements }: PasswordHintsProps) => {
  return (
    <div className={cx('flex', 'flex-col')}>
      <PasswordRequirement
        text={`Length is at least ${requirements.minLength} characters.`}
        isFulfilled={hasMinLength(value, requirements.minLength)}
      />
      <PasswordRequirement
        text={`Contains at least ${requirements.minLowercase} lower case character.`}
        isFulfilled={hasMinLowerCase(value, requirements.minLowercase)}
      />
      <PasswordRequirement
        text={`Contains at least ${requirements.minUppercase} upper case character.`}
        isFulfilled={hasMinUpperCase(value, requirements.minUppercase)}
      />
      <PasswordRequirement
        text={`Contains at least ${requirements.minNumbers} number.`}
        isFulfilled={hasMinNumbers(value, requirements.minNumbers)}
      />
      <PasswordRequirement
        text={`Contains at least ${requirements.minSymbols} symbol(s).`}
        isFulfilled={hasMinSymbols(value, requirements.minSymbols)}
      />
    </div>
  )
}

function hasMinLength(value: string, minimum: number): boolean {
  return value.length >= minimum
}

function hasMinLowerCase(value: string, minimum: number): boolean {
  const lowerCaseCount = Array.from(value).filter(
    (char) => char === char.toLowerCase() && char !== char.toUpperCase(),
  ).length
  return lowerCaseCount >= minimum
}

function hasMinUpperCase(value: string, minimum: number): boolean {
  const upperCaseCount = Array.from(value).filter(
    (char) => char === char.toUpperCase() && char !== char.toLowerCase(),
  ).length
  return upperCaseCount >= minimum
}

function hasMinNumbers(value: string, minimum: number): boolean {
  const numbersCount = Array.from(value).filter(
    (char) => !isNaN(Number(char)),
  ).length
  return numbersCount >= minimum
}

function hasMinSymbols(value: string, minimum: number): boolean {
  const symbolsCount = Array.from(value).filter(
    (char) => !char.match(/[a-zA-Z0-9\s]/),
  ).length
  return symbolsCount >= minimum
}

export type PasswordRequirementProps = {
  text: string
  isFulfilled?: boolean
}

export const PasswordRequirement = ({
  text,
  isFulfilled,
}: PasswordRequirementProps) => {
  return (
    <div
      className={cx(
        'flex flex-row',
        'gap-0.5',
        'items-center',
        { 'text-gray-400': !isFulfilled },
        { 'dark:text-gray-500': !isFulfilled },
        { 'text-green-500': isFulfilled },
        { 'dark:text-green-400': isFulfilled },
      )}
    >
      <IconCheck />
      <span>{text}</span>
    </div>
  )
}

export default PasswordHints
