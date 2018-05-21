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
  velocityPerManPerDay: number;
  parallels: number;
}

const homeInitialState: IHomeState = {
  isFetchingData: false,
  showManDayColumn: false,
  showTotalManDayColumn: false,
  showSPColumn: true,
  showTotalSPColumn: true,
  velocityPerManPerDay: 1,
  parallels: 2
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
    showTotalSPColumn: false,
    velocityPerManPerDay: 2 // FIXME
  }))
  .case(homeActionCreators.disableManDay, state => ({
    ...state,
    showManDayColumn: false,
    showTotalManDayColumn: false,
    showSPColumn: true,
    showTotalSPColumn: true,
    velocityPerManPerDay: 1
  }))
  .case(homeActionCreators.changeParallels, (state, payload) => ({
    ...state,
    parallels: payload
  }));

export const homeReducer = combineReducers({
  global: homeGlobalReducer,
  issueTable: issueTableReducer
});
