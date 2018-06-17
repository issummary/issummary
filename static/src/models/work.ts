import { Moment } from 'moment';
import { IMilestone } from './milestone';

export interface IWork {
  Issue: Issue;
  Label: ILabel;
  Dependencies: IDependencies;
  DependWorks: IWork[];
  TotalStoryPoint: number;
  StoryPoint: number;
}

export interface ILabel {
  ID: number;
  Name: string;
  Description: ILabelDescription;
  Parent: ILabel;
  Dependencies: IDependLabel[];
}

export interface ILabelDescription {
  Raw: string;
  DependLabelNames: string[];
  ParentName: string;
}

export interface Issue {
  ID: number;
  IID: number;
  DueDate?: Moment;
  Title: string;
  Description: IssueDescription;
  URL: string;
  ProjectName: string;
  Milestone: IMilestone;
}

export interface IssueDescription {
  Raw: string;
  DependencyIDs: IDependencyIDs;
  Summary: string;
  Note: string;
  Details: string;
}

export interface IDependencyIDs {
  issueIIDs: number[];
  labelNames: string[];
}

export interface IDependencies {
  Issues: Issue[];
  Labels: IDependLabel[];
}

export interface IDependLabel {
  Label: ILabel;
  RelatedIssues: Issue[];
}
