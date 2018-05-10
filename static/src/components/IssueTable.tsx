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
import { Work } from '../models/work';
import { CSSProperties } from 'react';

export interface IIssueTableProps {
  works: Work[];
  actions: IIssueTableActionCreators;
}

const rowStyle: CSSProperties = {
  wordWrap: 'break-word',
  whiteSpace: 'normal'
};

const IssueTableRow = (props: Work) => (
  <TableRow key={props.Issue.IID}>
    <TableRowColumn>{props.Issue.IID}</TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.Label && props.Label.Parent ? props.Label.Parent : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.Label ? props.Label : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>{props.Issue.Title}</TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.Issue.Summary ? props.Issue.Summary : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>{0}</TableRowColumn>
    <TableRowColumn style={rowStyle}>{'-'}</TableRowColumn>
  </TableRow>
);

class IssueTable extends React.Component<IIssueTableProps, undefined> {
  componentDidMount() {
    this.props.actions.requestUpdate();
  }

  render() {
    return (
      <Table fixedHeader={false} style={{ tableLayout: 'auto' }}>
        <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
          <TableRow>
            <TableHeaderColumn>IID</TableHeaderColumn>
            <TableHeaderColumn>Label</TableHeaderColumn>
            <TableHeaderColumn>Parent Label</TableHeaderColumn>
            <TableHeaderColumn>Title</TableHeaderColumn>
            <TableHeaderColumn>Summary</TableHeaderColumn>
            <TableHeaderColumn>SP</TableHeaderColumn>
            <TableHeaderColumn>Start Date</TableHeaderColumn>
          </TableRow>
        </TableHeader>
        <TableBody displayRowCheckbox={false}>
          {this.props.works.map(rowProp => IssueTableRow(rowProp))}
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
