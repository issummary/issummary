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

export interface IIssueTableProps {
  works: Work[];
  actions: IIssueTableActionCreators;
}

const IssueTableRow = (props: Work) => (
  <TableRow key={props.Issue.IID}>
    <TableRowColumn>{props.Issue.IID}</TableRowColumn>
    <TableRowColumn>{props.Issue.Title}</TableRowColumn>
    <TableRowColumn style={{ wordWrap: 'break-word', whiteSpace: 'normal' }}>
      {props.Issue.Summary}
    </TableRowColumn>
    <TableRowColumn style={{ wordWrap: 'break-word', whiteSpace: 'normal' }}>
      {props.Issue.Note}
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
