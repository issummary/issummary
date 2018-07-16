import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';

const actionCreator = actionCreatorFactory('BACKLOG');

export interface IBacklogActionCreators {
  enableManDay: ActionCreator<undefined>;
  disableManDay: ActionCreator<undefined>;
  changeParallels: ActionCreator<number>;
  changeProjectTextField: ActionCreator<string>;
}

export const backlogActionCreators: IBacklogActionCreators = {
  changeParallels: actionCreator<number>('CHANGE_PARALLELS'),
  changeProjectTextField: actionCreator<string>('CHANGE_PROJECT_TEXT_FIELD'),
  disableManDay: actionCreator<undefined>('DISABLE_MAN_DAY'),
  enableManDay: actionCreator<undefined>('ENABLE_MAN_DAY')
};
