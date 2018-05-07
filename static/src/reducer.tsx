import {combineReducers} from 'redux';
import {Action, combineActions, handleActions} from 'redux-actions';
import {ActionType, ICounterAmountPayload} from './actionCreators';

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

export interface IClasses {
    large: string;
    middle: string;
    small: string;
}

export interface IIssueTableRowProps {
    iid: number;
    classes: IClasses;
    title: string;
    description: string;
    summary: string;
    note: string;
}

export interface IIssueTableProps {
    rowProps: IIssueTableRowProps[];
}

const appInitialState: IAppState = { isOpenDrawer: false};
const counterInitialState: ICounterState = { count: 0 };
const issueTableInitialState: IIssueTableProps = { rowProps: [{
    iid: 100,
    classes: {
        large: "sss large",
        middle: "sss middle",
        small: "sss small",
    },
    title: "sample title",
    description: "sample description",
    summary: "sample summary",
    note: "sample note",
}] };

export const app = (state = appInitialState, action: Action<undefined>) => {
    const newState = Object.assign({}, state);
    switch (action.type) {
        case ActionType.TOGGLE_DRAWER:
            newState.isOpenDrawer = !newState.isOpenDrawer;
            return newState;
        default:
            return state;
    }
};

const combinedActions = combineActions(ActionType.INCREMENT, ActionType.DECREMENT);
export const counter = handleActions({
    [combinedActions](state: ICounterState, action: Action<ICounterAmountPayload>) {
        return (typeof action.payload === 'undefined') ? {...state} :
            { ...state, count: state.count + action.payload.amount };
    },
}, counterInitialState);

export const issueTable = (state = issueTableInitialState, action: Action<undefined>) => {
    const newState = Object.assign({}, state);
    return newState;
};

export const reducer = combineReducers({app, counter, issueTable});
