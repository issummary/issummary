import * as React from 'react';

import {
  Table,
  TableBody,
  TableHeader,
  TableHeaderColumn,
  TableRow,
  TableRowColumn
} from 'material-ui/Table';
import { IIssueTableActionCreators } from '../actions/issueTable';
import { CSSProperties } from 'react';
import { Dependencies, Issue, Work } from '../models/work';

export interface IIssueTableProps {
  works: Work[];
  actions: IIssueTableActionCreators;
}

export interface IIssueTableRowProps {
  work: Work;
  key: string;
}

const rowStyle: CSSProperties = {
  wordWrap: 'break-word',
  whiteSpace: 'normal'
};

const IssueIIDAndProjectName = (props: { issue: Issue }) => (
  <a href={props.issue.URL} target="_blank">
    {props.issue.ProjectName + ' #' + props.issue.IID}
  </a>
);

const IssueDependencies = (props: { deps: Dependencies }) => {
  // TODO: show label dependencies
  // const labels = props.deps.Labels;
  // const labelLinks = labels.map((l) => '~' + l.ID).join('->');

  const issues = props.deps.Issues;
  const issueLinks = issues.map(i => (
    <IssueIIDAndProjectName issue={i} key={i.ProjectName + i.IID} />
  ));

  if (issueLinks.length == 0) {
    return <span>-</span>;
  }

  const lastLink = issueLinks.pop();

  return (
    <span>
      {issueLinks.map((a, i) => (
        <span key={i}>
          {a}
          <span> </span>
        </span>
      ))}
      {lastLink}
    </span>
  );
};

const IssueTableRow = (props: IIssueTableRowProps) => (
  <TableRow key={props.key}>
    <TableRowColumn>
      <IssueIIDAndProjectName issue={props.work.Issue} />
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.work.Label && props.work.Label.Parent
        ? props.work.Label.Parent.Name
        : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.work.Label ? props.work.Label.Name : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>{props.work.Issue.Title}</TableRowColumn>
    <TableRowColumn style={{ ...rowStyle, width: 250 }}>
      {props.work.Issue.Description.Summary
        ? props.work.Issue.Description.Summary
        : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>{props.work.StoryPoint}</TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.work.Issue.DueDate
        ? props.work.Issue.DueDate.format('YYYY/MM/DD')
        : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>
      <IssueDependencies deps={props.work.Dependencies} />
    </TableRowColumn>
  </TableRow>
);

export class IssueTable extends React.Component<IIssueTableProps, any> {
  componentDidMount() {
    this.props.actions.requestUpdate();
  }

  render() {
    return (
      <Table fixedHeader={false} style={{ tableLayout: 'auto' }}>
        <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
          <TableRow>
            <TableHeaderColumn>Project+IID</TableHeaderColumn>
            <TableHeaderColumn>Parent Label</TableHeaderColumn>
            <TableHeaderColumn>Label</TableHeaderColumn>
            <TableHeaderColumn>Title</TableHeaderColumn>
            <TableHeaderColumn>Summary</TableHeaderColumn>
            <TableHeaderColumn>SP</TableHeaderColumn>
            <TableHeaderColumn>Due Date</TableHeaderColumn>
            <TableHeaderColumn>Deps</TableHeaderColumn>
          </TableRow>
        </TableHeader>
        <TableBody displayRowCheckbox={false}>
          {this.props.works.map(w => (
            <IssueTableRow work={w} key={w.Issue.ProjectName + w.Issue.IID} />
          ))}
        </TableBody>
      </Table>
    );
  }
}
