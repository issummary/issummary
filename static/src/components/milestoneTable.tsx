import * as React from 'react';

import { Table, TableBody } from 'material-ui/Table';
import TableHeader from 'material-ui/Table/TableHeader';
import { Milestone } from '../models/milestone';
import TableHeaderColumn from 'material-ui/Table/TableHeaderColumn';
import TableRowColumn from 'material-ui/Table/TableRowColumn';
import TableRow from 'material-ui/Table/TableRow';

export interface IMilestoneTableProps {
  milestones: Milestone[];
}

export class MilestoneTable extends React.Component<IMilestoneTableProps, any> {
  render() {
    console.log(this.props.milestones);
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
