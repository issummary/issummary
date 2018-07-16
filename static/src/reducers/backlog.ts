import { combineReducers } from 'redux';
import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { backlogActionCreators } from '../actions/backlog';
import { issueTableActionCreators, issueTableAsyncActionCreators } from '../actions/issueTable';
import { errorDialogReducer, IErrorDialogState } from './errorDialog';
import { IIssueTableState, issueTableReducer } from './issueTable';

export interface IBacklogState {
  isFetchingData: boolean;
  velocityPerManPerDay: number;
  parallels: number;
  selectedProjectName: string;
}

const backlogInitialState: IBacklogState = {
  isFetchingData: false,
  parallels: 2,
  selectedProjectName: 'All',
  velocityPerManPerDay: 1
};

const backlogReducer = reducerWithInitialState(backlogInitialState)
  .case(issueTableActionCreators.requestUpdate, state => ({
    ...state,
    isFetchingData: true
  }))
  .case(issueTableAsyncActionCreators.requestNewDataFetching.done, state => ({
    ...state,
    isFetchingData: false
  }))
  .case(issueTableAsyncActionCreators.requestNewDataFetching.failed, state => ({
    ...state,
    isFetchingData: false
  }))
  .case(backlogActionCreators.changeParallels, (state, payload) => ({
    ...state,
    parallels: payload
  }))
  .case(backlogActionCreators.changeProjectTextField, (state, payload) => ({
    ...state,
    selectedProjectName: payload
  }));

export interface IBacklogPageState {
  errorDialog: IErrorDialogState;
  global: IBacklogState;
  issueTable: IIssueTableState;
}

export const backlogPageReducer = combineReducers({
  errorDialog: errorDialogReducer,
  global: backlogReducer,
  issueTable: issueTableReducer
});
