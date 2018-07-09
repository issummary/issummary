import TableHeaderColumn from 'material-ui/Table/TableHeaderColumn';
import TableRow from 'material-ui/Table/TableRow';
import * as React from 'react';

interface IIssueTableHeadersProps {
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
  maxClassNum: number;
}

const generateClassColumns = (classNum: number) => {
  return Array.from(Array(classNum).keys()).map(n => (
    <TableHeaderColumn key={'class' + n}>Class{n + 1}</TableHeaderColumn>
  ));
};

// tslint:disable-next-line
export const IssueTableHeadersRow = (props: IIssueTableHeadersProps) => (
  <TableRow>
    <TableHeaderColumn>Project+IID</TableHeaderColumn>
    <TableHeaderColumn>Milestone</TableHeaderColumn>
    {generateClassColumns(props.maxClassNum)}
    <TableHeaderColumn>Title</TableHeaderColumn>
    <TableHeaderColumn>Summary</TableHeaderColumn>
    {props.showManDayColumn ? <TableHeaderColumn>ManDay</TableHeaderColumn> : null}
    {props.showTotalManDayColumn ? <TableHeaderColumn>Total MD</TableHeaderColumn> : null}
    {props.showTotalManDayColumn ? <TableHeaderColumn>Est. Date</TableHeaderColumn> : null}
    {props.showSPColumn ? <TableHeaderColumn>SP</TableHeaderColumn> : null}
    {props.showSPColumn ? <TableHeaderColumn>Total SP</TableHeaderColumn> : null}
    <TableHeaderColumn>Due Date</TableHeaderColumn>
    <TableHeaderColumn>Deps</TableHeaderColumn>
  </TableRow>
);
