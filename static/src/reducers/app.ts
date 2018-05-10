import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { appActionCreators } from '../actions/app';

export interface IAppState {
  isOpenDrawer: boolean;
}

const appInitialState: IAppState = { isOpenDrawer: false };

export const appReducer = reducerWithInitialState(appInitialState).case(
  appActionCreators.toggleDrawer,
  state => ({ ...state, isOpenDrawer: !state.isOpenDrawer })
);
