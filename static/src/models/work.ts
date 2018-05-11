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
  DueDate: string;
  Title: string;
  Description: IssueDescription;
  Summary: string;
  Note: string;
  URL: string;
}

export interface IssueDescription {
  Raw: string;
  DependencyIIDs: number[];
  Summary: string;
  Note: string;
  Details: string;
}

export interface Dependencies {
  Issues: Issue[];
  Labels: Label[];
}
