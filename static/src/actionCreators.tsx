import {ActionMap, createActions} from 'redux-actions';

export enum ActionType {
    ASYNC_INCREMENT = 'ASYNC_INCREMENT',
    INCREMENT       = 'INCREMENT',
    DECREMENT       = 'DECREMENT',
    TOGGLE_DRAWER   = 'TOGGLE_DRAWER',
}

export interface ICounterAmountPayload {
    amount: number;
}

const actions: ActionMap<ICounterAmountPayload | undefined, undefined> = {
    [ActionType.ASYNC_INCREMENT]: undefined,
    [ActionType.INCREMENT]      : (amount: number) => ({ amount: 1 }),
    [ActionType.DECREMENT]      : (amount: number) => ({ amount: -1 }),
    [ActionType.TOGGLE_DRAWER]  : undefined,
};

export const appActionCreator = createActions(actions);
