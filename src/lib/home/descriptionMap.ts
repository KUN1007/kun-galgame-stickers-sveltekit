export const descriptionMap: Record<string, string> = {
  '1-1': ''
}

// TODO:
for (let i = 1; i <= 6; i++) {
  for (let j = 1; j <= 80; j++) {
    descriptionMap[`${i}-${j}`] = ''
  }
}
