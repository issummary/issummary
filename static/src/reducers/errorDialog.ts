import { reducerWithInitialState } from 'typescript-fsa-reducers';
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
      error: payload.error,
      open: true,
    })
  )
  .case(errorDialogActionCreators.requestClosing, state => ({
    ...state,
    open: false
  }));
