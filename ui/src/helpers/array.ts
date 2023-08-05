export function stringArraysIdentical(a: string[], b: string[]): boolean {
  if (a.length !== b.length) {
    return false // Different lengths, not identical
  }

  // Sort the arrays to make sure the elements are in the same order
  const sortedArr1 = a.slice().sort()
  const sortedArr2 = b.slice().sort()

  // Compare each element in both arrays
  for (let i = 0; i < sortedArr1.length; i++) {
    if (sortedArr1[i] !== sortedArr2[i]) {
      return false // Different elements, not identical
    }
  }

  return true // Identical arrays
}
