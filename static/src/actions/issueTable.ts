import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';
import { Work } from '../models/work';

const actionCreator = actionCreatorFactory('ISSUE_TABLE');

export interface IIssueTableActionCreators {
  requestUpdate: ActionCreator<undefined>;
}

export const issueTableActionCreators: IIssueTableActionCreators = {
  requestUpdate: actionCreator<undefined>('REQUEST_UPDATE')
};

export const issueTableAsyncActionCreators = {
  requestNewDataFetching: actionCreator.async<{}, Work[], {}>(
    'REQUEST_NEW_DATA_FETCHING'
  )
};
