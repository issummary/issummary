import RaisedButton from 'material-ui/RaisedButton';
import * as React from 'react';
import { connect, Dispatch } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
  counterActionCreators,
  ICounterActionCreators
} from '../actionCreators';
import { IRootState } from '../reducer';
export interface ICounterProps {
  count: number;
  actions: ICounterActionCreators;
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
        <RaisedButton
          label="Async Increment"
          onClick={this.asyncIncrementClickEvent}
        />
        <RaisedButton label="Increment" onClick={this.incrementClickEvent} />
        <RaisedButton label="Decrement" onClick={this.decrementClickEvent} />
      </div>
    );
  }

  private asyncIncrementClickEvent(e: React.MouseEvent<{}>) {
    if (typeof this.props.actions !== 'undefined') {
      return this.props.actions.clickAsyncIncrementButton();
    }
  }

  private incrementClickEvent(e: React.MouseEvent<{}>) {
    if (typeof this.props.actions !== 'undefined') {
      return this.props.actions.clickIncrementButton();
    }
  }

  private decrementClickEvent(e: React.MouseEvent<{}>) {
    if (typeof this.props.actions !== 'undefined') {
      return this.props.actions.clickDecrementButton();
    }
  }
}

function mapStateToProps(state: IRootState) {
  return state.counter;
}

function mapDispatchToProps(dispatch: Dispatch<any>) {
  return { actions: bindActionCreators(counterActionCreators as {}, dispatch) };
}

// tslint:disable-next-line variable-name
export const ConnectedCounter = connect(mapStateToProps, mapDispatchToProps)(
  Counter as any
);
