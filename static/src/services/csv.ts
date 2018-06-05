import { Work } from '../models/work';
import * as moment from 'moment';
import koyomi = require('koyomi');
import { Moment } from 'moment';
import { eachSum } from './util';
import * as _ from 'lodash';

export const calcBizDay = (
  totalSP: number,
  velocityPerManPerDay: number,
  baseDay: Moment,
  parallels: number
): Date | undefined => {
  const totalManDay = totalSP / velocityPerManPerDay;
  const totalParallelManDay = Math.ceil(totalManDay / parallels);
  return koyomi.addBiz(baseDay.format('YYYY-MM-DD'), totalParallelManDay);
};

export const worksToCSV = (
  works: Work[],
  velocityPerManPerDay: number,
  baseDay: Moment,
  parallels: number
): string => {
  const totalSPs = eachSum(works.map(w => w.StoryPoint));

  // FIXME support story point column
  const header = [
    'ID',
    'Project',
    'IID',
    'Class1',
    'Class2',
    'Title',
    'Summary',
    'ManDay',
    'TotalDay',
    'EstDate',
    'DueDate',
    'IssueDepIIDs',
    'LabelDepIIDs'
  ];
  const lines = works.map((work, i) => {
    const bizRawDay = calcBizDay(
      totalSPs[i],
      velocityPerManPerDay,
      baseDay,
      parallels
    );
    const bizDayStr = bizRawDay
      ? moment(bizRawDay).format('YYYY-MM-DD')
      : '1年以上先';

    const labelIssues = _.flatMap(
      work.Dependencies.Labels,
      ls => ls.RelatedIssues
    );
    const uniqLabelIssues = _.uniqBy(labelIssues, i => i.ID);

    return [
      work.Issue.ID,
      work.Issue.ProjectName,
      work.Issue.IID,
      work.Label && work.Label.Parent ? work.Label.Parent : '-',
      work.Label ? work.Label : '-',
      work.Issue.Title,
      work.Issue.Description.Summary ? work.Issue.Description.Summary : '-',
      work.StoryPoint / velocityPerManPerDay,
      totalSPs[i] / velocityPerManPerDay,
      bizDayStr,
      work.Issue.DueDate ? work.Issue.DueDate : '-',
      work.Dependencies.Issues.map(i => `${i.ProjectName}#${i.IID}`).join('/'),
      uniqLabelIssues.map(i => `#${i.IID}`).join('/')
    ].join(',');
  });
  return header.join(',') + '\n' + lines.join('\n'); // FIXME ln
};
