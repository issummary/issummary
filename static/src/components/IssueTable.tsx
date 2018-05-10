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
import { CSSProperties } from 'react';
import { Work } from '../models/work';

export interface IIssueTableProps {
  works: Work[];
  actions: IIssueTableActionCreators;
}

export interface IIssueTableRowProps {
  work: Work;
  key: number;
}

const rowStyle: CSSProperties = {
  wordWrap: 'break-word',
  whiteSpace: 'normal'
};

const IssueTableRow = (props: IIssueTableRowProps) => (
  <TableRow>
    <TableRowColumn>{props.work.Issue.IID}</TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.work.Label && props.work.Label.Parent
        ? props.work.Label.Parent.Name
        : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.work.Label ? props.work.Label.Name : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>{props.work.Issue.Title}</TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.work.Issue.Summary ? props.work.Issue.Summary : '-'}
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
    console.log(this.props.works);
    return (
      <Table fixedHeader={false} style={{ tableLayout: 'auto' }}>
        <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
          <TableRow>
            <TableHeaderColumn>IID</TableHeaderColumn>
            <TableHeaderColumn>Parent Label</TableHeaderColumn>
            <TableHeaderColumn>Label</TableHeaderColumn>
            <TableHeaderColumn>Title</TableHeaderColumn>
            <TableHeaderColumn>Summary</TableHeaderColumn>
            <TableHeaderColumn>SP</TableHeaderColumn>
            <TableHeaderColumn>Start Date</TableHeaderColumn>
          </TableRow>
        </TableHeader>
        <TableBody displayRowCheckbox={false}>
          {this.props.works.map(w => (
            <IssueTableRow work={w} key={w.Issue.IID} />
          ))}
          {/*{this.props.works.map(rowProp => IssueTableRow(rowProp))}*/}
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
