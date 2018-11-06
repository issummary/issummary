import * as koyomi from 'koyomi';
import moment = require('moment');
import { IWork } from '../models/work';

export const sum = (arr: number[]): number => arr.reduce((a, b) => a + b, 0);
export const eachSum = (arr: number[]) => arr.map((e, i) => sum(arr.slice(0, i + 1)));

export const filterWorksByProjectNames = (works: IWork[], projectNames: string[]): IWork[] => {
  return works.filter(w => projectNames.find(p => p === w.Issue.ProjectName));
};

export const calcBizDayAsStr = (totalSP: number, velocityPerWeek: number, baseDay = moment()) => {
  const consumeWeeksNum = totalSP / velocityPerWeek;
  const consumeDaysNum = consumeWeeksNum * 5;
  const bizRawDayStr = koyomi.addBiz(baseDay.format('YYYY-MM-DD'), Math.ceil(consumeDaysNum));
  return bizRawDayStr ? moment(bizRawDayStr).format('YYYY-MM-DD') : '1年以上先';
};
