import * as _ from 'lodash';
import * as React from 'react';
import { Issue, IWork } from '../models/work';

// tslint:disable-next-line
export const IssueIIDAndProjectName = (props: { currentProjectName: string; issue: Issue }) => (
  <a href={props.issue.URL} target="_blank">
    {props.issue.ProjectName && props.issue.ProjectName !== props.currentProjectName
      ? props.issue.ProjectName + ' #' + props.issue.IID
      : '#' + props.issue.IID}
  </a>
);

// tslint:disable-next-line
const IssueDependencies = (props: { currentProjectName: string; issues: Issue[] }) => {
  const issueLinks = props.issues.map(i => (
    <IssueIIDAndProjectName currentProjectName={props.currentProjectName} issue={i} key={i.ProjectName + i.IID} />
  ));

  if (issueLinks.length === 0) {
    return null;
  }

  const lastLink = issueLinks.pop();

  return (
    <span>
      {issueLinks.map((a, i) => <span key={i}>{a} </span>)}
      {lastLink}
    </span>
  );
};

// tslint:disable-next-line
const LabelDependencies = (props: { currentProjectName: string; dependLabelName: string; dependIssues: Issue[] }) => {
  return (
    <span>
      {props.dependLabelName}(
      <IssueDependencies currentProjectName={props.currentProjectName} issues={props.dependIssues} />
      )
    </span>
  );
};

// tslint:disable-next-line
export const IssueTableIssueAndLabelDependenciesRow = (props: { work: IWork }) => {
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

  const labelWorks = labelOfIssueDescriptionDependWorks.concat(labelOfLabelDescriptionDependWorks);

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
        issues={issueOfIssueDescriptionDependWorks.map(iw => iw.Issue)}
      />
      {labelDependenciesDOMs.length > 0 ? ' ' : null}
      {labelDependenciesDOMs}
    </span>
  );
};
