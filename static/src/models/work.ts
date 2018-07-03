import { Moment } from 'moment';
import { IMilestone } from './milestone';

export interface IWork {
  Relation?: IWorkRelation;
  Issue: Issue;
  Label: ILabel;
  DependWorks: IWork[];
  TotalStoryPoint: number;
  StoryPoint: number;
}

export interface IWorkRelation {
  Type: string;
  LabelName: string;
}

export interface ILabel {
  ID: number;
  Name: string;
  Description: ILabelDescription;
  ParentName: string;
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
  Summary: string;
  Note: string;
  Details: string;
}
