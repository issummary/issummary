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
import { Milestone } from './models/milestone';
import { Work } from './models/work';
import { errorDialogActionCreators } from './actions/errorDialog';

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
  let works: Work[] = [];
  let milestones: Milestone[] = [];

  try {
    works = yield call(Api.fetchWorks);
  } catch (e) {
    yield put(issueTableAsyncActionCreators.requestNewDataFetching.failed(e));
    yield put(
      errorDialogActionCreators.failWorksResourceFetching({ error: e.Error })
    );
  }

  try {
    milestones = yield call(Api.fetchMilestones);
  } catch (e) {
    yield put(issueTableAsyncActionCreators.requestNewDataFetching.failed(e));
    yield put(
      errorDialogActionCreators.failWorksResourceFetching({ error: e.Error })
    );
  }

  return { milestones, works };
});

function* watchUpdateIssueTable() {
  yield takeEvery(issueTableActionCreators.requestUpdate.type, () =>
    requestNewIssueTableData(null)
  );
}

// single entry point to start all Sagas at once
export default function* rootSaga() {
  yield all([watchIncrementAsync(), watchUpdateIssueTable()]);
}
