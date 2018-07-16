import { combineReducers } from 'redux';
import { appReducer, IAppState } from './app';
import { counterReducer, ICounterState } from './counter';
import { IErrorDialogState } from './errorDialog';
import { homeReducer, IHomeState } from './home';
import { IIssueTableState } from './issueTable';

export interface ICombinedHomeState {
  // FIXME name
  errorDialog: IErrorDialogState;
  global: IHomeState;
  issueTable: IIssueTableState;
}

export interface IRootState {
  app: IAppState;
  counter: ICounterState;
  home: ICombinedHomeState;
}

export const reducer = combineReducers({
  app: appReducer,
  counter: counterReducer,
  home: homeReducer
});
