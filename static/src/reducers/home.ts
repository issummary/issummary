import { combineReducers } from 'redux';
import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { homeActionCreators } from '../actions/home';
import {
  issueTableActionCreators,
  issueTableAsyncActionCreators
} from '../actions/issueTable';
import { errorDialogReducer } from './errorDialog';
import { issueTableReducer } from './issueTable';

export interface IHomeState {
  isFetchingData: boolean;
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
  velocityPerManPerDay: number;
  parallels: number;
  selectedProjectName: string;
}

const homeInitialState: IHomeState = {
  isFetchingData: false,
  parallels: 2,
  selectedProjectName: 'All',
  showManDayColumn: false,
  showSPColumn: true,
  showTotalManDayColumn: false,
  showTotalSPColumn: true,
  velocityPerManPerDay: 1,
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
  .case(issueTableAsyncActionCreators.requestNewDataFetching.failed, state => ({
    ...state,
    isFetchingData: false
  }))
  .case(homeActionCreators.enableManDay, state => ({
    ...state,
    showManDayColumn: true,
    showSPColumn: false,
    showTotalManDayColumn: true,
    showTotalSPColumn: false,
    velocityPerManPerDay: 2 // FIXME
  }))
  .case(homeActionCreators.disableManDay, state => ({
    ...state,
    showManDayColumn: false,
    showSPColumn: true,
    showTotalManDayColumn: false,
    showTotalSPColumn: true,
    velocityPerManPerDay: 1
  }))
  .case(homeActionCreators.changeParallels, (state, payload) => ({
    ...state,
    parallels: payload
  }))
  .case(homeActionCreators.changeProjectTextField, (state, payload) => ({
    ...state,
    selectedProjectName: payload
  }));

export const homeReducer = combineReducers({
  errorDialog: errorDialogReducer,
  global: homeGlobalReducer,
  issueTable: issueTableReducer,
});
