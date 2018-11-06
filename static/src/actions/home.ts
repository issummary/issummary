import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';

const actionCreator = actionCreatorFactory('HOME');

export interface IHomeActionCreators {
  enableManDay: ActionCreator<undefined>;
  disableManDay: ActionCreator<undefined>;
  changeVelocityPerWeek: ActionCreator<number>;
  changeProjectTextField: ActionCreator<string>;
}

export const homeActionCreators: IHomeActionCreators = {
  changeProjectTextField: actionCreator<string>('CHANGE_PROJECT_TEXT_FIELD'),
  changeVelocityPerWeek: actionCreator<number>('CHANGE_VELOCITY_PER_WEEK'),
  disableManDay: actionCreator<undefined>('DISABLE_MAN_DAY'),
  enableManDay: actionCreator<undefined>('ENABLE_MAN_DAY')
};
