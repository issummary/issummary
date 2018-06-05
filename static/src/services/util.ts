import { Work } from '../models/work';

export const sum = (arr: number[]): number => arr.reduce((a, b) => a + b, 0);
export const eachSum = (arr: number[]) =>
  arr.map((e, i) => sum(arr.slice(0, i + 1)));

export const filterWorksByProjectNames = (
  works: Work[],
  projectNames: string[]
): Work[] => {
  return works.filter(w => projectNames.find(p => p === w.Issue.ProjectName));
};
