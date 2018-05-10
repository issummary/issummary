import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';
import { IIssueTableRowProps } from '../components/IssueTable';

export interface IIssueTablePayload {
  rows: IIssueTableRowProps;
}

const issueTableActionCreatorFactory = actionCreatorFactory('ISSUE_TABLE');
// TODO: issueTable用に書き換えていく

export interface IIssueTableActionCreators {
  clickAsyncIncrementButton: ActionCreator<undefined>;
  clickDecrementButton: ActionCreator<undefined>;
  clickIncrementButton: ActionCreator<undefined>;
  requestAmountChanging: ActionCreator<IIssueTablePayload>;
}

export const issueTableActionCreators: IIssueTableActionCreators = {
  clickAsyncIncrementButton: issueTableActionCreatorFactory<undefined>(
    'CLICK_ASYNC_INCREMENT_BUTTON'
  ),
  clickDecrementButton: issueTableActionCreatorFactory<undefined>(
    'CLICK_DECREMENT_BUTTON'
  ),
  clickIncrementButton: issueTableActionCreatorFactory<undefined>(
    'CLICK_INCREMENT_BUTTON'
  ),
  requestAmountChanging: issueTableActionCreatorFactory<IIssueTablePayload>(
    'REQUEST_AMOUNT_CHANGING'
  )
};

export const issueTableAsyncActionCreators = {
  changeAmountAsync: issueTableActionCreatorFactory.async<
    IIssueTablePayload,
    any,
    any
  >('CHANGE_AMOUNT_ASYNC')
};
