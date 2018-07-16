import { combineReducers } from 'redux';
import { appReducer, IAppState } from './app';
import { backlogPageReducer, IBacklogPageState } from './backlog';
import { counterReducer, ICounterState } from './counter';

export interface IRootState {
  app: IAppState;
  counter: ICounterState;
  backlogPage: IBacklogPageState;
}

export const reducer = combineReducers({
  app: appReducer,
  backlogPage: backlogPageReducer,
  counter: counterReducer
});
