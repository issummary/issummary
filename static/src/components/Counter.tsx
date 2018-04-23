import RaisedButton from 'material-ui/RaisedButton';
import * as React from 'react';
import {connect, Dispatch} from 'react-redux';
import {bindActionCreators} from 'redux';
import {Action, ActionFunction0, ActionFunctionAny} from 'redux-actions';
import {appActionCreator} from '../actionCreators';
import {ICounterState, IRootState} from '../reducer';
export interface ICounterProps {
    count?: number;
    actions?: {
        asyncIncrement(): ActionFunction0<Action<void>>;
        increment(): ActionFunction0<Action<void>>;
        decrement(): ActionFunction0<Action<void>>;
    };
}

export class Counter extends React.Component<ICounterProps, {}> {
    constructor(props: ICounterProps) {
        super(props);
        this.asyncIncrementClickEvent = this.asyncIncrementClickEvent.bind(this);
        this.incrementClickEvent = this.incrementClickEvent.bind(this);
        this.decrementClickEvent = this.decrementClickEvent.bind(this);
    }

    public render() {
        return (
            <div>
                <h1>Count: {this.props.count}</h1>
                <RaisedButton label='Async Increment' onClick={this.asyncIncrementClickEvent} />
                <RaisedButton label='Increment' onClick={this.incrementClickEvent} />
                <RaisedButton label='Decrement' onClick={this.decrementClickEvent} />
            </div>
        );
    }

    private asyncIncrementClickEvent(e: React.MouseEvent<{}>) {
        if (typeof(this.props.actions) !== 'undefined') {
            return this.props.actions.asyncIncrement();
        }
    }

    private incrementClickEvent(e: React.MouseEvent<{}>) {
        if (typeof(this.props.actions) !== 'undefined') {
            return this.props.actions.increment();
        }
    }

    private decrementClickEvent(e: React.MouseEvent<{}>) {
        if (typeof(this.props.actions) !== 'undefined') {
            return this.props.actions.decrement();
        }
    }
}

function mapStateToProps(state: IRootState) {
    return  state.counter;
}

interface IRootActionCreator {
    [actionName: string]: ActionFunctionAny<Action<undefined>>;
}

function mapDispatchToProps(dispatch: Dispatch<any>) {
    return { actions: bindActionCreators(appActionCreator, dispatch) };
}

// tslint:disable-next-line variable-name
export const ConnectedCounter = connect(mapStateToProps, mapDispatchToProps)(Counter as any);
