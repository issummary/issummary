import {
  IIssueTableProps,
  IIssueTableRowProps
} from '../components/IssueTable';
import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { issueTableActionCreators } from '../actions/issueTable';
import { Work } from '../models/work';

export interface IIssueTableState {
  rows: IIssueTableRowProps[];
}

const issueTableInitialState: IIssueTableState = {
  rows: [
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

const toTableRow = (work: Work): IIssueTableRowProps => ({
  iid: work.Issue.IID,
  classes: { large: 'large', middle: 'middle', small: 'small' },
  title: work.Issue.Title,
  description: work.Issue.Description,
  summary: work.Issue.Summary,
  note: work.Issue.Note
});

export const issueTableReducer = reducerWithInitialState(
  issueTableInitialState
).case(issueTableActionCreators.dataFetched, (state, payload) => ({
  rows: payload.map(toTableRow)
}));
