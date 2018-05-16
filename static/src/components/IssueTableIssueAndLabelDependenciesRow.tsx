import * as React from 'react';
import { Dependencies, DependLabel, Issue } from '../models/work';
import * as _ from 'lodash';

export const IssueIIDAndProjectName = (props: { issue: Issue }) => (
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

export const IssueTableIssueAndLabelDependenciesRow = (props: {
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
