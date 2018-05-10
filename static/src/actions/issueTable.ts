import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';
import { Work } from '../models/work';

const actionCreator = actionCreatorFactory('ISSUE_TABLE');

export interface IIssueTableActionCreators {
  dataFetched: ActionCreator<Work[]>;
  requestUpdate: ActionCreator<undefined>;
}

export const issueTableActionCreators: IIssueTableActionCreators = {
  dataFetched: actionCreator<Work[]>('DATA_FETCHED'),
  requestUpdate: actionCreator<undefined>('REQUEST_UPDATE')
};

export const issueTableAsyncActionCreators = {
  requestNewDataFetching: actionCreator.async<any, any, any>(
    'REQUEST_NEW_DATA_FETCHING'
  )
};
