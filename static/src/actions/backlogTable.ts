import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';
import { IMilestone } from '../models/milestone';
import { IWork } from '../models/work';

const actionCreator = actionCreatorFactory('ISSUE_TABLE');

export interface IBacklogTableActionCreators {
  requestUpdate: ActionCreator<undefined>;
}

export const backlogTableActionCreators: IBacklogTableActionCreators = {
  requestUpdate: actionCreator<undefined>('REQUEST_UPDATE')
};

export interface IRequestNewDataFetchingPayload {
  works: IWork[];
  milestones: IMilestone[];
}

export const backlogTableAsyncActionCreators = {
  requestNewDataFetching: actionCreator.async<null, IRequestNewDataFetchingPayload, null>('REQUEST_NEW_DATA_FETCHING') // FIXME
};
