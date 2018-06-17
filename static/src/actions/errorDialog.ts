import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';

const actionCreator = actionCreatorFactory('ERROR_DIALOG');

export interface IErrorDialogActionCreators {
  failWorksResourceFetching: ActionCreator<IFailWorksResourceFetchingPayload>;
  requestClosing: ActionCreator<undefined>;
}

export const errorDialogActionCreators: IErrorDialogActionCreators = {
  failWorksResourceFetching: actionCreator<IFailWorksResourceFetchingPayload>(
    'FAIL_WORKS_RESOURCE_FETCHING'
  ),
  requestClosing: actionCreator<undefined>('REQUEST_CLOSING')
};

export interface IFailWorksResourceFetchingPayload {
  error: string;
}
