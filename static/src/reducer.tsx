import {combineReducers} from 'redux';
import {Action, combineActions, handleActions} from 'redux-actions';
import {ActionType, ICounterAmountPayload} from './actionCreators';

export interface IRootState {
    app: IAppState;
    counter: ICounterState;
}

export interface IAppState {
    isOpenDrawer: boolean;
}

export interface ICounterState {
    count: number;
}

const appInitialState: IAppState = { isOpenDrawer: false};
const counterInitialState: ICounterState = { count: 0 };

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

export const reducer = combineReducers({app, counter});
