import * as React from 'react';
import * as _ from 'lodash';
import { CSSProperties } from 'react';
import { Dependencies, DependLabel, Issue, Work } from '../models/work';
import TableRowColumn from 'material-ui/Table/TableRowColumn';
import TableRow from 'material-ui/Table/TableRow';

export interface IIssueTableRowProps {
  work: Work;
  key: string;
  totalSP: number;
}

const rowStyle: CSSProperties = {
  wordWrap: 'break-word',
  whiteSpace: 'normal'
};

const IssueIIDAndProjectName = (props: { issue: Issue }) => (
  <a href={props.issue.URL} target="_blank">
    {props.issue.ProjectName
      ? props.issue.ProjectName + ' #' + props.issue.IID
      : '#' + props.issue.IID}
  </a>
);

const IssueDependencies = (props: { issues: Issue[] }) => {
  const issueLinks = props.issues.map(i => (
    <IssueIIDAndProjectName issue={i} key={i.ProjectName + i.IID} />
  ));

  if (issueLinks.length == 0) {
    return null;
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

const LabelDependencies = (props: { dependLabel: DependLabel }) => {
  return (
    <span>
      {props.dependLabel.Label.Name}(
      <IssueDependencies issues={props.dependLabel.RelatedIssues} />
      )
    </span>
  );
};

const IssueAndLabelDependencies = (props: {
  deps: Dependencies;
  labelDeps: DependLabel[];
}) => {
  if (
    props.deps.Issues.length == 0 &&
    props.deps.Labels.length == 0 &&
    props.labelDeps.length == 0
  ) {
    return <span>-</span>;
  }

  const labels = props.deps.Labels.concat(props.labelDeps);
  const uniqueLabels = _.uniqBy(labels, l => l.Label.Name);

  return (
    <span>
      <IssueDependencies issues={props.deps.Issues} />
      {uniqueLabels.map(l => (
        <LabelDependencies
          dependLabel={l}
          key={'LabelDependencies' + l.Label.ID}
        />
      ))}
    </span>
  );
};

export const IssueTableRow = (props: IIssueTableRowProps) => (
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
    <TableRowColumn style={rowStyle}>{props.totalSP}</TableRowColumn>
    <TableRowColumn style={rowStyle}>
      {props.work.Issue.DueDate
        ? props.work.Issue.DueDate.format('YYYY/MM/DD')
        : '-'}
    </TableRowColumn>
    <TableRowColumn style={rowStyle}>
      <IssueAndLabelDependencies
        deps={props.work.Dependencies}
        labelDeps={props.work.Label ? props.work.Label.Dependencies : []}
      />
    </TableRowColumn>
  </TableRow>
);
