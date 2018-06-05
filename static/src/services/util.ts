export const sum = (arr: number[]): number => arr.reduce((a, b) => a + b, 0);
export const eachSum = (arr: number[]) =>
  arr.map((e, i) => sum(arr.slice(0, i + 1)));
