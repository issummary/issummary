import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';

const actionCreator = actionCreatorFactory('COUNTER');

export interface ICounterAmountPayload {
  amount: number;
}

export interface ICounterActionCreators {
  clickAsyncIncrementButton: ActionCreator<undefined>;
  clickDecrementButton: ActionCreator<undefined>;
  clickIncrementButton: ActionCreator<undefined>;
  requestAmountChanging: ActionCreator<ICounterAmountPayload>;
}

export const counterActionCreators: ICounterActionCreators = {
  clickAsyncIncrementButton: actionCreator<undefined>(
    'CLICK_ASYNC_INCREMENT_BUTTON'
  ),
  clickDecrementButton: actionCreator<undefined>('CLICK_DECREMENT_BUTTON'),
  clickIncrementButton: actionCreator<undefined>('CLICK_INCREMENT_BUTTON'),
  requestAmountChanging: actionCreator<ICounterAmountPayload>(
    'REQUEST_AMOUNT_CHANGING'
  )
};

export const counterAsyncActionCreators = {
  changeAmountAsync: actionCreator.async<ICounterAmountPayload, any, any>(
    'CHANGE_AMOUNT_ASYNC'
  )
};
