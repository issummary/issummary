import { combineReducers } from 'redux';
import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { backlogActionCreators } from '../actions/backlog';
import { backlogTableActionCreators, backlogTableAsyncActionCreators } from '../actions/backlogTable';
import { backlogTableReducer, IBacklogTableState } from './backlogTable';
import { errorDialogReducer, IErrorDialogState } from './errorDialog';

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
  .case(backlogTableActionCreators.requestUpdate, state => ({
    ...state,
    isFetchingData: true
  }))
  .case(backlogTableAsyncActionCreators.requestNewDataFetching.done, state => ({
    ...state,
    isFetchingData: false
  }))
  .case(backlogTableAsyncActionCreators.requestNewDataFetching.failed, state => ({
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
  backlogTable: IBacklogTableState;
}

export const backlogPageReducer = combineReducers({
  backlogTable: backlogTableReducer,
  errorDialog: errorDialogReducer,
  global: backlogReducer
});
