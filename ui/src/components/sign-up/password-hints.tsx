import { IconCheck } from '@koupr/ui'
import cx from 'classnames'
import { AccountPasswordRequirements } from '@/client/idp/account'

export type PasswordHintsProps = {
  value: string
  requirements: AccountPasswordRequirements
}

const PasswordHints = ({ value, requirements }: PasswordHintsProps) => {
  return (
    <div className={cx('flex', 'flex-col')}>
      <PasswordRequirement
        text={`Length is at least ${requirements.minLength} characters.`}
        isFulfilled={hasMinLength(value, requirements.minLength)}
      />
      <PasswordRequirement
        text={`Contains at least ${requirements.minLowercase} lowercase character.`}
        isFulfilled={hasMinLowerCase(value, requirements.minLowercase)}
      />
      <PasswordRequirement
        text={`Contains at least ${requirements.minUppercase} uppercase character.`}
        isFulfilled={hasMinUpperCase(value, requirements.minUppercase)}
      />
      <PasswordRequirement
        text={`Contains at least ${requirements.minNumbers} number.`}
        isFulfilled={hasMinNumbers(value, requirements.minNumbers)}
      />
      <PasswordRequirement
        text={`Contains at least ${requirements.minSymbols} special character(s) (!#$%).`}
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
    (character) =>
      character === character.toLowerCase() &&
      character !== character.toUpperCase(),
  ).length
  return lowerCaseCount >= minimum
}

function hasMinUpperCase(value: string, minimum: number): boolean {
  const upperCaseCount = Array.from(value).filter(
    (character) =>
      character === character.toUpperCase() &&
      character !== character.toLowerCase(),
  ).length
  return upperCaseCount >= minimum
}

function hasMinNumbers(value: string, minimum: number): boolean {
  const numbersCount = Array.from(value).filter(
    (character) => !isNaN(Number(character)),
  ).length
  return numbersCount >= minimum
}

function hasMinSymbols(value: string, minimum: number): boolean {
  const symbolsCount = Array.from(value).filter(
    (character) => !character.match(/[a-zA-Z0-9\s]/),
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
