import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { issueTableAsyncActionCreators } from '../actions/issueTable';
import { Work } from '../models/work';

export interface IIssueTableState {
  works: Work[];
}

const issueTableInitialState: IIssueTableState = { works: [] };

export const issueTableReducer = reducerWithInitialState(
  issueTableInitialState
).case(
  issueTableAsyncActionCreators.requestNewDataFetching.done,
  (state, payload) => ({ works: payload.result })
);
