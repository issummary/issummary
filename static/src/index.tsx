import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { AppContainer } from 'react-hot-loader';

import { Provider } from 'react-redux';
import { applyMiddleware, compose, createStore } from 'redux';
import createSagaMiddleware from 'redux-saga';

import App from './components/App';
import { reducer } from './reducer';
import rootSaga from './sagas';

const sagaMiddleware = createSagaMiddleware();
const composeEnhancers =
  (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const store = createStore(
  reducer,
  composeEnhancers(applyMiddleware(sagaMiddleware))
);

sagaMiddleware.run(rootSaga);

ReactDOM.render(
  <AppContainer>
    <MuiThemeProvider>
      <Provider store={store}>
        <App />
      </Provider>
    </MuiThemeProvider>
  </AppContainer>,
  document.getElementById('example')
);

if (module.hot) {
  module.hot.accept('./components/App', () => {
    const NextApp = (require('./components/App') as any).default; // tslint:disable-line variable-name
    ReactDOM.render(
      <AppContainer>
        <MuiThemeProvider>
          <Provider store={store}>
            <NextApp />
          </Provider>
        </MuiThemeProvider>
      </AppContainer>,
      document.getElementById('example')
    );
  });
}
