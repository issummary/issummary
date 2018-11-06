import TableRow from 'material-ui/Table/TableRow';
import TableRowColumn from 'material-ui/Table/TableRowColumn';
import * as React from 'react';
import { CSSProperties } from 'react';
import { IWork } from '../models/work';
import { calcBizDayAsStr } from '../services/util';
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
  velocityPerWeek: number;
  maxClassNum: number;
}

const rowStyle: CSSProperties = {
  whiteSpace: 'normal',
  wordWrap: 'break-word'
};

// tslint:disable-next-line
const TotalSPPoint = (props: { work: IWork; velocityPerManPerDay: number }) => (
  <span>
    <br />(+{props.work.TotalStoryPoint / props.velocityPerManPerDay})
  </span>
);

// tslint:disable-next-line
export const IssueTableRow = (props: IIssueTableRowProps) => {
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
      <TableRowColumn style={rowStyle}>{calcBizDayAsStr(props.totalSP, props.velocityPerWeek)}</TableRowColumn>
      <TableRowColumn style={rowStyle}>
        {props.work.Issue.DueDate ? props.work.Issue.DueDate.format('YYYY/MM/DD') : '-'}
      </TableRowColumn>
      <TableRowColumn style={rowStyle}>
        <IssueTableIssueAndLabelDependenciesRow work={props.work} />
      </TableRowColumn>
    </TableRow>
  );
};
