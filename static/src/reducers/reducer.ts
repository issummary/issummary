import { combineReducers } from 'redux';
import { IIssueTableProps } from '../components/IssueTable';
import { appReducer, IAppState } from './app';
import { counterReducer, ICounterState } from './counter';
import { issueTableReducer } from './issueTable';

export interface IRootState {
  app: IAppState;
  counter: ICounterState;
  issueTable: IIssueTableProps;
}

export const reducer = combineReducers({
  app: appReducer,
  counter: counterReducer,
  issueTable: issueTableReducer
});
