import { combineReducers } from 'redux';
import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { appActionCreators, counterActionCreators } from './actionCreators';
import {
  IClasses,
  IIssueTableProps,
  IIssueTableRowProps
} from './components/IssueTable';

export interface IRootState {
  app: IAppState;
  counter: ICounterState;
  issueTable: IIssueTableProps;
}

export interface IAppState {
  isOpenDrawer: boolean;
}

export interface ICounterState {
  count: number;
}

export type IIssueTableState = IIssueTableProps;

const appInitialState: IAppState = { isOpenDrawer: false };
const counterInitialState: ICounterState = { count: 0 };
const issueTableInitialState: IIssueTableState = {
  rowProps: [
    {
      iid: 1,
      classes: { large: 'large', middle: 'middle', small: 'small' },
      title: 'sample title',
      description: 'description sample',
      summary: 'sample summary',
      note: 'sample note'
    }
  ]
};

const app = reducerWithInitialState(appInitialState).case(
  appActionCreators.toggleDrawer,
  state => ({ ...state, isOpenDrawer: !state.isOpenDrawer })
);

const counter = reducerWithInitialState(counterInitialState)
  .case(counterActionCreators.clickIncrementButton, state => ({
    ...state,
    count: state.count + 1
  }))
  .case(counterActionCreators.clickDecrementButton, state => ({
    ...state,
    count: state.count - 1
  }))
  .case(counterActionCreators.requestAmountChanging, (state, payload) => ({
    ...state,
    count: state.count + payload.amount
  }));

const issueTable = reducerWithInitialState(issueTableInitialState);

export const reducer = combineReducers({ app, counter, issueTable });
