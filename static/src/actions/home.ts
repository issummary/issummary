import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';

const actionCreator = actionCreatorFactory('HOME');

export interface IHomeActionCreators {
  enableManDay: ActionCreator<undefined>;
  disableManDay: ActionCreator<undefined>;
  changeParallels: ActionCreator<number>;
  changeProjectTextField: ActionCreator<string>;
}

export const homeActionCreators: IHomeActionCreators = {
  enableManDay: actionCreator<undefined>('ENABLE_MAN_DAY'),
  disableManDay: actionCreator<undefined>('DISABLE_MAN_DAY'),
  changeParallels: actionCreator<number>('CHANGE_PARALLELS'),
  changeProjectTextField: actionCreator<string>('CHANGE_PROJECT_TEXT_FIELD')
};
