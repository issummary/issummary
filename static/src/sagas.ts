import { delay, SagaIterator } from 'redux-saga';
import { all, call, put, takeEvery } from 'redux-saga/effects';
import { Action } from 'typescript-fsa';
import { bindAsyncAction } from 'typescript-fsa-redux-saga';
import {
  counterActionCreators,
  counterAsyncActionCreators,
  ICounterAmountPayload
} from './actions/counter';
import {
  issueTableActionCreators,
  issueTableAsyncActionCreators
} from './actions/issueTable';
import { Api } from './services/api';

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

const requestNewIssueTableData = bindAsyncAction(
  issueTableAsyncActionCreators.requestNewDataFetching
)(function*(): SagaIterator {
  return yield call(Api.fetchWorks);
});

function* watchUpdateIssueTable() {
  yield takeEvery(
    issueTableActionCreators.requestUpdate.type,
    (a: Action<ICounterAmountPayload>) => requestNewIssueTableData(a.payload)
  );
}

// single entry point to start all Sagas at once
export default function* rootSaga() {
  yield all([watchIncrementAsync(), watchUpdateIssueTable()]);
}
