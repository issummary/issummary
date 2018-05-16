import * as React from 'react';
import TableRow from 'material-ui/Table/TableRow';
import TableHeaderColumn from 'material-ui/Table/TableHeaderColumn';

interface IIssueTableHeadersProps {
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
}

export const IssueTableHeadersRow = (props: IIssueTableHeadersProps) => (
  <TableRow>
    <TableHeaderColumn>Project+IID</TableHeaderColumn>
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
    {props.showSPColumn ? <TableHeaderColumn>SP</TableHeaderColumn> : null}
    {props.showSPColumn ? (
      <TableHeaderColumn>Total SP</TableHeaderColumn>
    ) : null}
    <TableHeaderColumn>Due Date</TableHeaderColumn>
    <TableHeaderColumn>Deps</TableHeaderColumn>
  </TableRow>
);
