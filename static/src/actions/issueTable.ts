import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';
import { Work } from '../models/work';
import { Milestone } from '../models/milestone';

const actionCreator = actionCreatorFactory('ISSUE_TABLE');

export interface IIssueTableActionCreators {
  requestUpdate: ActionCreator<undefined>;
}

export const issueTableActionCreators: IIssueTableActionCreators = {
  requestUpdate: actionCreator<undefined>('REQUEST_UPDATE')
};

export interface RequestNewDataFetchingPayload {
  works: Work[];
  milestones: Milestone[];
}

export const issueTableAsyncActionCreators = {
  requestNewDataFetching: actionCreator.async<
    null,
    RequestNewDataFetchingPayload,
    null
  >('REQUEST_NEW_DATA_FETCHING') // FIXME
};
