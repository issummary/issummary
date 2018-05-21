import * as React from 'react';
import { CSSProperties } from 'react';
import { Work } from '../models/work';
import TableRowColumn from 'material-ui/Table/TableRowColumn';
import TableRow from 'material-ui/Table/TableRow';
import {
  IssueIIDAndProjectName,
  IssueTableIssueAndLabelDependenciesRow
} from './IssueTableIssueAndLabelDependenciesRow';
import koyomi = require('koyomi');
import moment from 'moment';
export interface IIssueTableRowProps {
  work: Work;
  key: string;
  totalSP: number;
  velocityPerManPerDay: number;
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
}

const rowStyle: CSSProperties = {
  wordWrap: 'break-word',
  whiteSpace: 'normal'
};

const today = moment().format('YYYY-MM-DD');

export const IssueTableRow = (props: IIssueTableRowProps) => {
  const totalManDay = props.totalSP / props.velocityPerManPerDay;
  const bizRawDay = koyomi.addBiz(today, Math.ceil(totalManDay));

  const bizDay = bizRawDay
    ? moment(bizRawDay).format('YYYY-MM-DD')
    : '1年以上先';
  return (
    <TableRow key={props.key}>
      <TableRowColumn>
        <IssueIIDAndProjectName issue={props.work.Issue} />
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        {props.work.Label && props.work.Label.Parent
          ? props.work.Label.Parent.Name
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
        <IssueTableIssueAndLabelDependenciesRow
          deps={props.work.Dependencies}
          labelDeps={props.work.Label ? props.work.Label.Dependencies : []}
        />
      </TableRowColumn>
    </TableRow>
  );
};
