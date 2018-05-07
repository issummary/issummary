import {ActionMap, createActions} from 'redux-actions';

export enum ActionType {
    ASYNC_INCREMENT = 'ASYNC_INCREMENT',
    INCREMENT       = 'INCREMENT',
    DECREMENT       = 'DECREMENT',
    TOGGLE_DRAWER   = 'TOGGLE_DRAWER',
    UPDATE_TABLE   = 'UPDATE_TABLE',
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

const issueTableActions: ActionMap<ICounterAmountPayload | undefined, undefined> = {
};

export const issueTableActionCreator = createActions(issueTableActions);
