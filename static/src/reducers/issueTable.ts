import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { issueTableActionCreators } from '../actions/issueTable';
import { Work } from '../models/work';

export interface IIssueTableState {
  works: Work[];
}

const issueTableInitialState: IIssueTableState = { works: [] };

export const issueTableReducer = reducerWithInitialState(
  issueTableInitialState
).case(issueTableActionCreators.dataFetched, (state, works) => ({ works }));
