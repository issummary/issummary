import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { issueTableAsyncActionCreators } from '../actions/issueTable';
import { IMilestone } from '../models/milestone';
import { IWork } from '../models/work';

export interface IIssueTableState {
  works: IWork[];
  milestones: IMilestone[];
}

const issueTableInitialState: IIssueTableState = { works: [], milestones: [] };

export const issueTableReducer = reducerWithInitialState(
  issueTableInitialState
).case(
  issueTableAsyncActionCreators.requestNewDataFetching.done,
  (state, payload) => ({
    milestones: payload.result.milestones,
    works: payload.result.works,
  })
);
