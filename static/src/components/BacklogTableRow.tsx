import koyomi = require('koyomi');
import TableRow from 'material-ui/Table/TableRow';
import TableRowColumn from 'material-ui/Table/TableRowColumn';
import * as moment from 'moment';
import * as React from 'react';
import { CSSProperties } from 'react';
import { IWork } from '../models/work';
import { BacklogTableIssueAndLabelDependenciesRow } from './BacklogTableIssueAndLabelDependenciesRow';
export interface IBacklogTableRowProps {
  work: IWork;
  key: string;
  totalSP: number;
  velocityPerManPerDay: number;
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
  parallels: number;
  maxClassNum: number;
}

const rowStyle: CSSProperties = {
  whiteSpace: 'normal',
  wordWrap: 'break-word'
};

const today = moment().format('YYYY-MM-DD');

// tslint:disable-next-line
const TotalSPPoint = (props: { work: IWork; velocityPerManPerDay: number }) => (
  <span>
    <br />(+{props.work.TotalStoryPoint / props.velocityPerManPerDay})
  </span>
);

// tslint:disable-next-line
export const BacklogTableRow = (props: IBacklogTableRowProps) => {
  const totalManDay = props.totalSP / props.velocityPerManPerDay;
  const totalParallelManDay = Math.ceil(totalManDay / props.parallels);
  const bizRawDay = koyomi.addBiz(today, totalParallelManDay);

  const bizDay = bizRawDay ? moment(bizRawDay).format('YYYY-MM-DD') : '1年以上先';

  const dashList: string[] = new Array<string>(props.maxClassNum).fill('-');
  const label = props.work.Label;
  let classLabelNames: string[] = label ? label.ParentNames.concat([label.Name]) : [];
  classLabelNames = classLabelNames.concat(dashList).slice(0, props.maxClassNum);

  // tslint:disable-next-line
  const ClassColumns = classLabelNames.map(name => (
    <TableRowColumn key={'class' + name} style={rowStyle}>
      {name}
    </TableRowColumn>
  ));

  return (
    <TableRow key={props.key}>
      <TableRowColumn>
        <a href={props.work.Issue.URL} target="_blank">
          {props.work.Issue.ProjectName ? props.work.Issue.ProjectName : null}
          #{props.work.Issue.IID}
        </a>
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        {props.work.Issue.Milestone ? props.work.Issue.Milestone.Title : '-'}
      </TableRowColumn>
      {ClassColumns}
      <TableRowColumn style={rowStyle}>{props.work.Issue.Title}</TableRowColumn>
      <TableRowColumn style={{ ...rowStyle, width: 250 }}>
        {props.work.Issue.Description.Summary ? props.work.Issue.Description.Summary : '-'}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        {props.work.StoryPoint / props.velocityPerManPerDay}
        {props.work.TotalStoryPoint !== 0 ? (
          <TotalSPPoint work={props.work} velocityPerManPerDay={props.velocityPerManPerDay} />
        ) : null}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>{props.totalSP / props.velocityPerManPerDay}</TableRowColumn>
      {props.showTotalManDayColumn ? <TableRowColumn style={rowStyle}>{bizDay}</TableRowColumn> : null}
      <TableRowColumn style={rowStyle}>
        {props.work.Issue.DueDate ? props.work.Issue.DueDate.format('YYYY/MM/DD') : '-'}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        <BacklogTableIssueAndLabelDependenciesRow work={props.work} />
      </TableRowColumn>
    </TableRow>
  );
};
