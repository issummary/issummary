import { combineReducers } from 'redux';
import { IIssueTableProps } from '../components/IssueTable';
import { appReducer, IAppState } from './app';
import { counterReducer, ICounterState } from './counter';
import { issueTableReducer } from './issueTable';
import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { issueTableActionCreators } from '../actions/issueTable';

export interface IRootState {
  app: IAppState;
  counter: ICounterState;
  home: {
    issueTable: IIssueTableProps;
  };
}

export interface IHomeState {
  isFetchingData: boolean;
}

const homeInitialState: IHomeState = { isFetchingData: false };

const homeGlobalReducer = reducerWithInitialState(homeInitialState)
  .case(issueTableActionCreators.requestUpdate, state => ({
    ...state,
    isFetchingData: true
  }))
  .case(issueTableActionCreators.dataFetched, state => ({
    ...state,
    isFetchingData: false
  }));

const homeReducer = combineReducers({
  global: homeGlobalReducer,
  issueTable: issueTableReducer
});

export const reducer = combineReducers({
  app: appReducer,
  counter: counterReducer,
  home: homeReducer
});
