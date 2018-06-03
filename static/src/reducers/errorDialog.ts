import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { issueTableAsyncActionCreators } from '../actions/issueTable';
import { Work } from '../models/work';
import { Milestone } from '../models/milestone';
import { errorDialogActionCreators } from '../actions/errorDialog';

export interface IErrorDialogState {
  open: boolean;
  error: string;
}

const errorDialogInitialState: IErrorDialogState = { open: false, error: '' };

export const errorDialogReducer = reducerWithInitialState(
  errorDialogInitialState
)
  .case(
    errorDialogActionCreators.failWorksResourceFetching,
    (state, payload) => ({
      ...state,
      open: true,
      error: payload.error
    })
  )
  .case(errorDialogActionCreators.requestClosing, state => ({
    ...state,
    open: false
  }));
