import * as _ from 'lodash';
import * as React from 'react';
import { Issue, IWork } from '../models/work';

export const IssueIIDAndProjectName = (props: {
  // tslint:disable-line
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

const IssueDependencies = (props: {
  // tslint:disable-line
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
          <span> </span>
        </span>
      ))}
      {lastLink}
    </span>
  );
};

const LabelDependencies = (props: {
  // tslint:disable-line
  currentProjectName: string;
  dependLabelName: string;
  dependIssues: Issue[];
}) => {
  return (
    <span>
      {props.dependLabelName}(
      <IssueDependencies
        currentProjectName={props.currentProjectName}
        issues={props.dependIssues}
      />
      )
    </span>
  );
};

export const IssueTableIssueAndLabelDependenciesRow = (props: {
  // tslint:disable-line
  work: IWork;
}) => {
  const dependWorks = props.work.DependWorks;

  const issueOfIssueDescriptionDependWorks = dependWorks.filter(
    w => w.Relation && w.Relation.Type === 'IssueOfIssueDescription'
  ); // FIXME
  const labelOfIssueDescriptionDependWorks = dependWorks.filter(
    w => w.Relation && w.Relation.Type === 'LabelOfIssueDescription'
  ); // FIXME
  const labelOfLabelDescriptionDependWorks = dependWorks.filter(
    w => w.Relation && w.Relation.Type === 'LabelOfLabelDescription'
  ); // FIXME

  if (
    issueOfIssueDescriptionDependWorks.length === 0 &&
    labelOfIssueDescriptionDependWorks.length === 0 &&
    labelOfLabelDescriptionDependWorks.length === 0
  ) {
    return <span>-</span>;
  }

  const labelWorks = labelOfIssueDescriptionDependWorks.concat(
    labelOfLabelDescriptionDependWorks
  );

  const groupedWorks = _.groupBy(labelWorks, w => w.Relation!.LabelName);

  const labelDependenciesDOMs = Object.keys(groupedWorks).map(labelName => {
    const works = groupedWorks[labelName];
    const firstWork = works[0];
    return (
      <LabelDependencies
        currentProjectName={firstWork.Issue.ProjectName} // FIXME
        dependLabelName={labelName}
        dependIssues={works.map(lw => lw.Issue)}
        key={'LabelDependencies' + labelName}
      />
    );
  });

  return (
    <span>
      <IssueDependencies
        currentProjectName={props.work.Issue.ProjectName}
        issues={props.work.DependWorks.map(dw => dw.Issue)}
      />
      {labelDependenciesDOMs}
    </span>
  );
};
