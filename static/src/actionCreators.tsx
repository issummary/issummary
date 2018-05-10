import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';

const counterActionCreatorFactory = actionCreatorFactory('COUNTER');

export interface IAppActionCreators {
  toggleDrawer: ActionCreator<undefined>;
}

export const appActionCreators = {
  toggleDrawer: counterActionCreatorFactory('TOGGLE_DRAWER')
};

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
  clickAsyncIncrementButton: counterActionCreatorFactory<undefined>(
    'CLICK_ASYNC_INCREMENT_BUTTON'
  ),
  clickDecrementButton: counterActionCreatorFactory<undefined>(
    'CLICK_DECREMENT_BUTTON'
  ),
  clickIncrementButton: counterActionCreatorFactory<undefined>(
    'CLICK_INCREMENT_BUTTON'
  ),
  requestAmountChanging: counterActionCreatorFactory<ICounterAmountPayload>(
    'REQUEST_AMOUNT_CHANGING'
  )
};

export const counterAsyncActionCreators = {
  changeAmountAsync: counterActionCreatorFactory.async<
    ICounterAmountPayload,
    any,
    any
  >('CHANGE_AMOUNT_ASYNC')
};
