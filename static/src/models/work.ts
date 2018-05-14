import { Moment } from 'moment';

export interface Work {
  Issue: Issue;
  Label: Label;
  Dependencies: Dependencies;
  StoryPoint: number;
}

export interface Label {
  ID: number;
  Name: string;
  Description: LabelDescription;
  Parent: Label;
  Dependencies: Label[];
}

export interface LabelDescription {
  Raw: string;
  DependLabelNames: string[];
  ParentName: string;
}

export interface Issue {
  ID: number;
  IID: number;
  DueDate: Moment;
  Title: string;
  Description: IssueDescription;
  URL: string;
  ProjectName: string;
}

export interface IssueDescription {
  Raw: string;
  DependencyIDs: DependencyIDs;
  Summary: string;
  Note: string;
  Details: string;
}

export interface DependencyIDs {
  issueIIDs: number[];
  labelNames: string[];
}

export interface Dependencies {
  Issues: Issue[];
  Labels: Label[];
}
