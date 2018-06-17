import * as _ from 'lodash';
import * as React from 'react';
import { IDependencies, IDependLabel, Issue } from '../models/work';

export const IssueIIDAndProjectName = (props: {// tslint:disable-line
  currentProjectName: string;
  issue: Issue;
}) => (
  <a href={props.issue.URL} target="_blank">
    {props.issue.ProjectName &&
    props.issue.ProjectName !== props.currentProjectName
      ? props.issue.ProjectName + ' #' + props.issue.IID
      : '#' + props.issue.IID}
  </a>
);


const IssueDependencies = (props: { // tslint:disable-line
  currentProjectName: string;
  issues: Issue[];
}) => {
  const issueLinks = props.issues.map(i => (
    <IssueIIDAndProjectName
      currentProjectName={props.currentProjectName}
      issue={i}
      key={i.ProjectName + i.IID}
    />
  ));

  if (issueLinks.length === 0) {
    return null;
  }

  const lastLink = issueLinks.pop();

  return (
    <span>
      {issueLinks.map((a, i) => (
        <span key={i}>
          {a}
          <span>{' '}</span>
        </span>
      ))}
      {lastLink}
    </span>
  );
};

const LabelDependencies = (props: {// tslint:disable-line
  currentProjectName: string;
  dependLabel: IDependLabel;
}) => {
  return (
    <span>
      {props.dependLabel.Label.Name}(
      <IssueDependencies
        currentProjectName={props.currentProjectName}
        issues={props.dependLabel.RelatedIssues}
      />
      )
    </span>
  );
};

export const IssueTableIssueAndLabelDependenciesRow = (props: {// tslint:disable-line
  currentProjectName: string;
  deps: IDependencies;
  labelDeps: IDependLabel[];
}) => {
  if (
    props.deps.Issues.length === 0 &&
    props.deps.Labels.length === 0 &&
    props.labelDeps.length === 0
  ) {
    return <span>-</span>;
  }

  const labels = props.deps.Labels.concat(props.labelDeps);
  const uniqueLabels = _.uniqBy(labels, l => l.Label.Name);

  return (
    <span>
      <IssueDependencies
        currentProjectName={props.currentProjectName}
        issues={props.deps.Issues}
      />
      {uniqueLabels.map(l => (
        <LabelDependencies
          currentProjectName={props.currentProjectName}
          dependLabel={l}
          key={'LabelDependencies' + l.Label.ID}
        />
      ))}
    </span>
  );
};
