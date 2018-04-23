import React = require('react');
import ReactDOM = require('react-dom');
import injectTapEventPlugin = require('react-tap-event-plugin');

(window as any).ReactDOM = ReactDOM;
(window as any).React = React;

injectTapEventPlugin();
