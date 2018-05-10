import { delay, SagaIterator } from 'redux-saga';
import { all, call, put, takeEvery } from 'redux-saga/effects';
import { Action } from 'typescript-fsa';
import { bindAsyncAction } from 'typescript-fsa-redux-saga';
import {
  counterActionCreators,
  counterAsyncActionCreators,
  ICounterAmountPayload
} from './actionCreators';

function* incrementAsync(payload: ICounterAmountPayload) {
  yield delay(1000);
  yield put(counterActionCreators.requestAmountChanging(payload));
}

const counterIncrementWorker = bindAsyncAction(
  counterAsyncActionCreators.changeAmountAsync
)(function*(payload: ICounterAmountPayload): SagaIterator {
  yield call(incrementAsync, { ...payload, amount: 1 });
});

function* watchIncrementAsync() {
  yield takeEvery(
    counterActionCreators.clickAsyncIncrementButton.type,
    (a: Action<ICounterAmountPayload>) => counterIncrementWorker(a.payload)
  );
}

// single entry point to start all Sagas at once
export default function* rootSaga() {
  yield all([watchIncrementAsync()]);
}
