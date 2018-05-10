import * as React from 'react';

import {
  Table,
  TableBody,
  TableHeader,
  TableHeaderColumn,
  TableRow,
  TableRowColumn
} from 'material-ui/Table';
import { connect, Dispatch } from 'react-redux';
import { IRootState } from '../reducers/reducer';
import { bindActionCreators } from 'redux';
import {
  IIssueTableActionCreators,
  issueTableActionCreators
} from '../actions/issueTable';

export interface IClasses {
  large: string;
  middle: string;
  small: string;
}

export interface IIssueTableRowProps {
  iid: number;
  classes: IClasses;
  title: string;
  description: string;
  summary: string;
  note: string;
}

export interface IIssueTableProps {
  rows: IIssueTableRowProps[];
  actions: IIssueTableActionCreators;
}

const IssueTableRow = (props: IIssueTableRowProps) => (
  <TableRow key={props.iid}>
    <TableRowColumn>{props.iid}</TableRowColumn>
    <TableRowColumn>{props.classes.large}</TableRowColumn>
    <TableRowColumn>{props.classes.middle}</TableRowColumn>
    <TableRowColumn>{props.classes.small}</TableRowColumn>
    <TableRowColumn>{props.title}</TableRowColumn>
    <TableRowColumn style={{ wordWrap: 'break-word', whiteSpace: 'normal' }}>
      {props.summary}
    </TableRowColumn>
    <TableRowColumn style={{ wordWrap: 'break-word', whiteSpace: 'normal' }}>
      {props.note}
    </TableRowColumn>
  </TableRow>
);

class IssueTable extends React.Component<IIssueTableProps, undefined> {
  componentDidMount() {
    this.props.actions.requestUpdate();
  }

  render() {
    return (
      <Table>
        <TableHeader>
          <TableRow>
            <TableHeaderColumn>IID</TableHeaderColumn>
            <TableHeaderColumn>Large Class</TableHeaderColumn>
            <TableHeaderColumn>Middle Class</TableHeaderColumn>
            <TableHeaderColumn>Small Class</TableHeaderColumn>
            <TableHeaderColumn>Title</TableHeaderColumn>
            <TableHeaderColumn>Summary</TableHeaderColumn>
            <TableHeaderColumn>Note</TableHeaderColumn>
          </TableRow>
        </TableHeader>
        <TableBody>
          {this.props.rows.map(rowProp => IssueTableRow(rowProp))}
        </TableBody>
      </Table>
    );
  }
}

function mapStateToProps(state: IRootState) {
  return state.issueTable;
}

function mapDispatchToProps(dispatch: Dispatch<any>) {
  return {
    actions: bindActionCreators(issueTableActionCreators as {}, dispatch)
  };
}

// tslint:disable-next-line variable-name
export const ConnectedIssueTable = connect(mapStateToProps, mapDispatchToProps)(
  IssueTable as any
);
