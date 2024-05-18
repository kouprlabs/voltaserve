export function isMicrosoftOfficeLockFile(name: string) {
  return name.startsWith('~$')
}

export function isOpenOfficeOfficeLockFile(name: string) {
  return name.startsWith('.~lock.') && name.endsWith('#')
}
