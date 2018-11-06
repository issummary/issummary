import * as _ from 'lodash';
import { Moment } from 'moment';
import { IWork } from '../models/work';
import { calcBizDayAsStr, eachSum } from './util';

export const worksToCSV = (
  works: IWork[],
  velocityPerManPerDay: number,
  baseDay: Moment,
  velocityPerWeek: number
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
    const bizDayStr = calcBizDayAsStr(totalSPs[index], velocityPerWeek, baseDay);
    const labelIssues = work.DependWorks.filter(w => w.Relation && w.Relation.Type === 'LabelOfLabelDescription').map(
      w => w.Issue
    );
    const dependIssues = work.DependWorks.filter(w => w.Relation && w.Relation.Type === 'IssueOfIssueDescription').map(
      w => w.Issue
    );

    const uniqLabelIssues = _.uniqBy(labelIssues, i => i.ID);

    return [
      work.Issue.ID,
      work.Issue.ProjectName,
      work.Issue.IID,
      work.Label && work.Label.Description.ParentName ? work.Label.Description.ParentName : '-',
      work.Label ? work.Label.Name : '-',
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
