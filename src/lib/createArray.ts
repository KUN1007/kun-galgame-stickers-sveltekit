export const createArray = (n: number): number[] => {
  return Array.from({ length: n }, (_, index) => index + 1)
}
