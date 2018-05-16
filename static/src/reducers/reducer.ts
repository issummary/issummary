import { combineReducers } from 'redux';
import { IIssueTableProps } from '../components/IssueTable';
import { appReducer, IAppState } from './app';
import { counterReducer, ICounterState } from './counter';
import { homeReducer } from './home';

export interface IRootState {
  app: IAppState;
  counter: ICounterState;
  home: {
    issueTable: IIssueTableProps;
  };
}

export const reducer = combineReducers({
  app: appReducer,
  counter: counterReducer,
  home: homeReducer
});
