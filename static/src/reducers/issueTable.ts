import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { backlogActionCreators } from '../actions/backlog';
import { issueTableAsyncActionCreators } from '../actions/issueTable';
import { IMilestone } from '../models/milestone';
import { IWork } from '../models/work';

export interface IIssueTableState {
  works: IWork[];
  milestones: IMilestone[];
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
}

const issueTableInitialState: IIssueTableState = {
  milestones: [],
  showManDayColumn: false,
  showSPColumn: true,
  showTotalManDayColumn: false,
  showTotalSPColumn: true,
  works: []
};

export const issueTableReducer = reducerWithInitialState(issueTableInitialState)
  .case(issueTableAsyncActionCreators.requestNewDataFetching.done, (state, payload) => ({
    ...state,
    milestones: payload.result.milestones,
    works: payload.result.works
  }))
  .case(backlogActionCreators.enableManDay, state => ({
    ...state,
    showManDayColumn: true,
    showSPColumn: false,
    showTotalManDayColumn: true,
    showTotalSPColumn: false,
    velocityPerManPerDay: 2 // FIXME
  }))
  .case(backlogActionCreators.disableManDay, state => ({
    ...state,
    showManDayColumn: false,
    showSPColumn: true,
    showTotalManDayColumn: false,
    showTotalSPColumn: true,
    velocityPerManPerDay: 1
  }));
