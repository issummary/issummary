import { IIssueTableProps } from '../components/IssueTable';
import { reducerWithInitialState } from 'typescript-fsa-reducers';

export type IIssueTableState = IIssueTableProps;

const issueTableInitialState: IIssueTableState = {
  rowProps: [
    {
      iid: 1,
      classes: { large: 'large', middle: 'middle', small: 'small' },
      title: 'sample title',
      description: 'description sample',
      summary: 'sample summary',
      note: 'sample note'
    }
  ]
};

export const issueTableReducer = reducerWithInitialState(
  issueTableInitialState
);
