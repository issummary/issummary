import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';
import { IMilestone } from '../models/milestone';
import { IWork } from '../models/work';

const actionCreator = actionCreatorFactory('ISSUE_TABLE');

export interface IIssueTableActionCreators {
  requestUpdate: ActionCreator<undefined>;
}

export const issueTableActionCreators: IIssueTableActionCreators = {
  requestUpdate: actionCreator<undefined>('REQUEST_UPDATE')
};

export interface IRequestNewDataFetchingPayload {
  works: IWork[];
  milestones: IMilestone[];
}

export const issueTableAsyncActionCreators = {
  requestNewDataFetching: actionCreator.async<
    null,
    IRequestNewDataFetchingPayload,
    null
  >('REQUEST_NEW_DATA_FETCHING') // FIXME
};
