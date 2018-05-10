import actionCreatorFactory, { ActionCreator } from 'typescript-fsa';

const actionCreator = actionCreatorFactory('APP');

export interface IAppActionCreators {
  toggleDrawer: ActionCreator<undefined>;
}

export const appActionCreators = {
  toggleDrawer: actionCreator('TOGGLE_DRAWER')
};
