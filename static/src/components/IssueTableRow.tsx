import koyomi = require('koyomi');
import TableRow from 'material-ui/Table/TableRow';
import TableRowColumn from 'material-ui/Table/TableRowColumn';
import * as moment from 'moment';
import * as React from 'react';
import { CSSProperties } from 'react';
import { IWork } from '../models/work';
import { IssueTableIssueAndLabelDependenciesRow } from './IssueTableIssueAndLabelDependenciesRow';
export interface IIssueTableRowProps {
  work: IWork;
  key: string;
  totalSP: number;
  velocityPerManPerDay: number;
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
  parallels: number;
}

const rowStyle: CSSProperties = {
  whiteSpace: 'normal',
  wordWrap: 'break-word'
};

const today = moment().format('YYYY-MM-DD');

const TotalSPPoint = (
  props: { work: IWork; velocityPerManPerDay: number } // tslint:disable-line
) => (
  <span>
    <br />(+{props.work.TotalStoryPoint / props.velocityPerManPerDay})
  </span>
);

export const IssueTableRow = (props: IIssueTableRowProps) => {
  // tslint:disable-line
  const totalManDay = props.totalSP / props.velocityPerManPerDay;
  const totalParallelManDay = Math.ceil(totalManDay / props.parallels);
  const bizRawDay = koyomi.addBiz(today, totalParallelManDay);

  const bizDay = bizRawDay
    ? moment(bizRawDay).format('YYYY-MM-DD')
    : '1年以上先';

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
      <TableRowColumn style={rowStyle}>
        {props.work.Label && props.work.Label.Description.ParentName
          ? props.work.Label.Description.ParentName
          : '-'}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        {props.work.Label ? props.work.Label.Name : '-'}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>{props.work.Issue.Title}</TableRowColumn>
      <TableRowColumn style={{ ...rowStyle, width: 250 }}>
        {props.work.Issue.Description.Summary
          ? props.work.Issue.Description.Summary
          : '-'}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        {props.work.StoryPoint / props.velocityPerManPerDay}
        {props.work.TotalStoryPoint !== 0 ? (
          <TotalSPPoint
            work={props.work}
            velocityPerManPerDay={props.velocityPerManPerDay}
          />
        ) : null}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        {props.totalSP / props.velocityPerManPerDay}
      </TableRowColumn>
      {props.showTotalManDayColumn ? (
        <TableRowColumn style={rowStyle}>{bizDay}</TableRowColumn>
      ) : null}
      <TableRowColumn style={rowStyle}>
        {props.work.Issue.DueDate
          ? props.work.Issue.DueDate.format('YYYY/MM/DD')
          : '-'}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        <IssueTableIssueAndLabelDependenciesRow work={props.work} />
      </TableRowColumn>
    </TableRow>
  );
};
