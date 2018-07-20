import * as _ from 'lodash';
import * as React from 'react';
import { branch, compose, mapProps, renderComponent } from 'recompose';
import { Issue, IWork } from '../models/work';

interface IBaseProjectNameAndIssueNumberProps {
  projectName: string;
  issueURL: string;
  issueNumber: number;
}

// tslint:disable-next-line
const BaseProjectNameAndIssueNumber = ({ projectName, issueURL, issueNumber }: IBaseProjectNameAndIssueNumberProps) => (
  <a href={issueURL} target="_blank">
    {projectName}#{issueNumber}
  </a>
);

interface IProjectNameAndIssueNumberProps {
  currentProjectName: string;
  issue: Issue;
}

// tslint:disable-next-line
const ProjectNameAndIssueNumber = compose<IBaseProjectNameAndIssueNumberProps, IProjectNameAndIssueNumberProps>(
  mapProps(({ currentProjectName, issue }: IProjectNameAndIssueNumberProps) => ({
    issueNumber: issue.IID,
    issueProjectName: issue.ProjectName,
    issueURL: issue.URL,
    projectName: issue.ProjectName !== currentProjectName ? issue.ProjectName : null
  }))
)(BaseProjectNameAndIssueNumber);

interface IIssueDependencies {
  currentProjectName: string;
  issues: Issue[];
}

interface IBaseIssueDependencies extends IIssueDependencies {
  lastIssue: Issue;
}

// tslint:disable-next-line
const BaseIssueDependencies = (props: IBaseIssueDependencies) => {
  const issueLinks = props.issues.map(i => (
    <ProjectNameAndIssueNumber currentProjectName={props.currentProjectName} issue={i} />
  ));

  return (
    <span>
      {issueLinks.map((a, i) => <span key={i}>{a} </span>)}
      <ProjectNameAndIssueNumber currentProjectName={props.currentProjectName} issue={props.lastIssue} />
    </span>
  );
};

// tslint:disable-next-line
const IssueDependencies = mapProps(({ currentProjectName, issues }: IIssueDependencies) => ({
  currentProjectName,
  issues: issues.slice(0, issues.length - 1),
  lastIssue: issues[issues.length - 1]
}))(BaseIssueDependencies);

interface ILabelDependenciesProps {
  currentProjectName: string;
  dependLabelName: string;
  dependIssues: Issue[];
}

// tslint:disable-next-line
const LabelDependencies = (props: ILabelDependenciesProps) => {
  return (
    <span>
      {props.dependLabelName}(
      <IssueDependencies currentProjectName={props.currentProjectName} issues={props.dependIssues} />
      )
    </span>
  );
};

interface IBaseBacklogTableIssueAndLabelDependenciesRowProps {
  currentProjectName: string;
  issueOfIssueDescriptionDependWorks: IWork[];
  labelOfIssueDescriptionDependWorks: IWork[];
  labelOfLabelDescriptionDependWorks: IWork[];
}

// tslint:disable-next-line
const BaseBacklogTableIssueAndLabelDependenciesRow = (props: IBaseBacklogTableIssueAndLabelDependenciesRowProps) => {
  const labelWorks = props.labelOfIssueDescriptionDependWorks.concat(props.labelOfLabelDescriptionDependWorks);
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
        currentProjectName={props.currentProjectName}
        issues={props.issueOfIssueDescriptionDependWorks.map(iw => iw.Issue)}
      />
      {labelDependenciesDOMs.length > 0 ? ' ' : null}
      {labelDependenciesDOMs}
    </span>
  );
};

interface IBacklogTableIssueAndLabelDependenciesRowProps {
  work: IWork;
}

// tslint:disable-next-line
export const BacklogTableIssueAndLabelDependenciesRow = compose<
  IBaseBacklogTableIssueAndLabelDependenciesRowProps,
  IBacklogTableIssueAndLabelDependenciesRowProps
>(
  mapProps(({ work }: IBacklogTableIssueAndLabelDependenciesRowProps) => {
    const dependWorks = work.DependWorks;

    const issueOfIssueDescriptionDependWorks = dependWorks.filter(
      w => w.Relation && w.Relation.Type === 'IssueOfIssueDescription'
    );
    const labelOfIssueDescriptionDependWorks = dependWorks.filter(
      w => w.Relation && w.Relation.Type === 'LabelOfIssueDescription'
    );
    const labelOfLabelDescriptionDependWorks = dependWorks.filter(
      w => w.Relation && w.Relation.Type === 'LabelOfLabelDescription'
    );

    return {
      currentProjectName: work.Issue.ProjectName,
      issueOfIssueDescriptionDependWorks,
      labelOfIssueDescriptionDependWorks,
      labelOfLabelDescriptionDependWorks
    };
  }),
  branch(
    (
      props: IBaseBacklogTableIssueAndLabelDependenciesRowProps // FIXME extract to method
    ) =>
      props.issueOfIssueDescriptionDependWorks.length === 0 &&
      props.labelOfIssueDescriptionDependWorks.length === 0 &&
      props.labelOfLabelDescriptionDependWorks.length === 0,
    renderComponent(() => <span>-</span>)
  )
)(BaseBacklogTableIssueAndLabelDependenciesRow);
