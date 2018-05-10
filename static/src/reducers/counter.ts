import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { counterActionCreators } from '../actions/counter';

export interface ICounterState {
  count: number;
}

const counterInitialState: ICounterState = { count: 0 };

export const counterReducer = reducerWithInitialState(counterInitialState)
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
