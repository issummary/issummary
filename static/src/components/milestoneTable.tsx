import * as React from 'react';

import { Table, TableBody } from 'material-ui/Table';
import TableHeader from 'material-ui/Table/TableHeader';
import TableHeaderColumn from 'material-ui/Table/TableHeaderColumn';
import TableRow from 'material-ui/Table/TableRow';
import TableRowColumn from 'material-ui/Table/TableRowColumn';
import { IMilestone } from '../models/milestone';

export interface IMilestoneTableProps {
  milestones: IMilestone[];
}

export class MilestoneTable extends React.Component<IMilestoneTableProps, any> {
  public render() {
    console.log(this.props.milestones);// tslint:disable-line
    return (
      <Table fixedHeader={false} style={{ tableLayout: 'auto' }}>
        <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
          <TableRow>
            <TableHeaderColumn>IID</TableHeaderColumn>
            <TableHeaderColumn>Title</TableHeaderColumn>
            <TableHeaderColumn>Start Date</TableHeaderColumn>
            <TableHeaderColumn>Due Date</TableHeaderColumn>
          </TableRow>
        </TableHeader>
        <TableBody displayRowCheckbox={false}>
          {this.props.milestones.map(m => {
            return (
              <TableRow key={m.IID}>
                <TableRowColumn>{m.IID}</TableRowColumn>
                <TableRowColumn>{m.Title}</TableRowColumn>
                <TableRowColumn>
                  {m.StartDate ? m.StartDate.format('YYYY/MM/DD') : '-'}
                </TableRowColumn>
                <TableRowColumn>
                  {m.DueDate ? m.DueDate.format('YYYY/MM/DD') : '-'}
                </TableRowColumn>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    );
  }
}
