import {delay} from 'redux-saga';
import {all, put, takeEvery} from 'redux-saga/effects';
import {ActionType, appActionCreator} from './actionCreators';

function* watchAsyncIncrement(): any {
    yield delay(1000);
    yield put(appActionCreator.increment());
}

export function* watchIncrementAsync() {
    yield takeEvery(ActionType.ASYNC_INCREMENT, watchAsyncIncrement);
}

// single entry point to start all Sagas at once
export default function* rootSaga() {
    yield all([
        watchIncrementAsync(),
    ]);
}
