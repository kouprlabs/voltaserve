export function stringArraysIdentical(a: string[], b: string[]): boolean {
  if (a.length !== b.length) {
    return false // Different lengths, not identical
  }

  // Sort the arrays to make sure the elements are in the same order
  const sortedA = a.slice().sort()
  const sortedB = b.slice().sort()

  // Compare each element in both arrays
  for (let i = 0; i < sortedA.length; i++) {
    if (sortedA[i] !== sortedB[i]) {
      return false // Different elements, not identical
    }
  }

  return true // Identical arrays
}
