import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';
import { Work } from '../models/work';
import { Milestone } from '../models/milestone';

const actionCreator = actionCreatorFactory('ERROR_DIALOG');

export interface IErrorDialogActionCreators {
  failWorksResourceFetching: ActionCreator<failWorksResourceFetchingPayload>;
  requestClosing: ActionCreator<undefined>;
}

export const errorDialogActionCreators: IErrorDialogActionCreators = {
  failWorksResourceFetching: actionCreator<failWorksResourceFetchingPayload>(
    'FAIL_WORKS_RESOURCE_FETCHING'
  ),
  requestClosing: actionCreator<undefined>('REQUEST_CLOSING')
};

export interface failWorksResourceFetchingPayload {
  error: string;
}
