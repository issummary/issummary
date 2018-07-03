import koyomi = require('koyomi');
import * as _ from 'lodash';
import * as moment from 'moment';
import { Moment } from 'moment';
import { IWork } from '../models/work';
import { eachSum } from './util';

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
  works: IWork[],
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
  const lines = works.map((work, index) => {
    const bizRawDay = calcBizDay(
      totalSPs[index],
      velocityPerManPerDay,
      baseDay,
      parallels
    );
    const bizDayStr = bizRawDay
      ? moment(bizRawDay).format('YYYY-MM-DD')
      : '1年以上先';

    const labelIssues = work.DependWorks.filter(
      w => w.Relation && w.Relation.Type === 'LabelOfLabelDescription'
    ).map(w => w.Issue);
    const dependIssues = work.DependWorks.filter(
      w => w.Relation && w.Relation.Type === 'IssueOfIssueDescription'
    ).map(w => w.Issue);

    const uniqLabelIssues = _.uniqBy(labelIssues, i => i.ID);

    return [
      work.Issue.ID,
      work.Issue.ProjectName,
      work.Issue.IID,
      work.Label && work.Label.ParentName ? work.Label.ParentName : '-',
      work.Label ? work.Label : '-',
      work.Issue.Title,
      work.Issue.Description.Summary ? work.Issue.Description.Summary : '-',
      work.StoryPoint / velocityPerManPerDay,
      totalSPs[index] / velocityPerManPerDay,
      bizDayStr,
      work.Issue.DueDate ? work.Issue.DueDate : '-',
      dependIssues.map(i => `${i.ProjectName}#${i.IID}`).join('/'),
      uniqLabelIssues.map(i => `#${i.IID}`).join('/')
    ].join(',');
  });
  return header.join(',') + '\n' + lines.join('\n'); // FIXME ln
};
