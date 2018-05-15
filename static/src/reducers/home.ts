import {
  issueTableActionCreators,
  issueTableAsyncActionCreators
} from '../actions/issueTable';
import { combineReducers } from 'redux';
import { issueTableReducer } from './issueTable';
import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { homeActionCreators } from '../actions/home';

export interface IHomeState {
  isFetchingData: boolean;
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
}

const homeInitialState: IHomeState = {
  isFetchingData: false,
  showManDayColumn: false,
  showTotalManDayColumn: false,
  showSPColumn: true,
  showTotalSPColumn: true
};

const homeGlobalReducer = reducerWithInitialState(homeInitialState)
  .case(issueTableActionCreators.requestUpdate, state => ({
    ...state,
    isFetchingData: true
  }))
  .case(issueTableAsyncActionCreators.requestNewDataFetching.done, state => ({
    ...state,
    isFetchingData: false
  }))
  .case(homeActionCreators.enableManDay, state => ({
    ...state,
    showManDayColumn: true,
    showTotalManDayColumn: true,
    showSPColumn: false,
    showTotalSPColumn: false
  }))
  .case(homeActionCreators.disableManDay, state => ({
    ...state,
    showManDayColumn: false,
    showTotalManDayColumn: false,
    showSPColumn: true,
    showTotalSPColumn: true
  }));

export const homeReducer = combineReducers({
  global: homeGlobalReducer,
  issueTable: issueTableReducer
});
