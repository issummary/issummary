import TableHeaderColumn from 'material-ui/Table/TableHeaderColumn';
import TableRow from 'material-ui/Table/TableRow';
import * as React from 'react';

interface IIssueTableHeadersProps {
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
}

export const IssueTableHeadersRow = (props: IIssueTableHeadersProps) => (// tslint:disable-line
  <TableRow>
    <TableHeaderColumn>Project+IID</TableHeaderColumn>
    <TableHeaderColumn>Milestone</TableHeaderColumn>
    <TableHeaderColumn>Parent Label</TableHeaderColumn>
    <TableHeaderColumn>Label</TableHeaderColumn>
    <TableHeaderColumn>Title</TableHeaderColumn>
    <TableHeaderColumn>Summary</TableHeaderColumn>
    {props.showManDayColumn ? (
      <TableHeaderColumn>ManDay</TableHeaderColumn>
    ) : null}
    {props.showTotalManDayColumn ? (
      <TableHeaderColumn>Total MD</TableHeaderColumn>
    ) : null}
    {props.showTotalManDayColumn ? (
      <TableHeaderColumn>Est. Date</TableHeaderColumn>
    ) : null}
    {props.showSPColumn ? <TableHeaderColumn>SP</TableHeaderColumn> : null}
    {props.showSPColumn ? (
      <TableHeaderColumn>Total SP</TableHeaderColumn>
    ) : null}
    <TableHeaderColumn>Due Date</TableHeaderColumn>
    <TableHeaderColumn>Deps</TableHeaderColumn>
  </TableRow>
);
