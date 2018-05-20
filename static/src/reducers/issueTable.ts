import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { issueTableAsyncActionCreators } from '../actions/issueTable';
import { Work } from '../models/work';
import { Milestone } from '../models/milestone';

export interface IIssueTableState {
  works: Work[];
  milestones: Milestone[];
}

const issueTableInitialState: IIssueTableState = { works: [], milestones: [] };

export const issueTableReducer = reducerWithInitialState(
  issueTableInitialState
).case(
  issueTableAsyncActionCreators.requestNewDataFetching.done,
  (state, payload) => ({
    works: payload.result.works,
    milestones: payload.result.milestones
  })
);
