import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import MenuItem from 'material-ui/MenuItem';
import * as React from 'react';
import { connect, Dispatch } from 'react-redux';
import { BrowserRouter as Router, Link, Route } from 'react-router-dom';
import { bindActionCreators } from 'redux';
import { IRootState } from '../reducers/reducer';
import { About } from './About';
import { ConnectedCounter } from './Counter';
import { Home } from './Home';
import { appActionCreators, IAppActionCreators } from '../actions/app';

export interface IAppProps {
  isOpenDrawer: boolean;
  actions: {
    toggleDrawer(): IAppActionCreators;
  };
}

class App extends React.Component<IAppProps, undefined> {
  constructor(props: IAppProps) {
    super(props);
    this.handleLeftIconButtonTouchTap = this.handleLeftIconButtonTouchTap.bind(
      this
    );
    this.handleClose = this.handleClose.bind(this);
    this.handleRequestChange = this.handleRequestChange.bind(this);
  }

  public render() {
    return (
      <Router>
        <div>
          <AppBar
            title="Issummary"
            iconClassNameRight="muidocs-icon-navigation-expand-more"
            onLeftIconButtonClick={this.handleLeftIconButtonTouchTap}
          />

          <Drawer
            docked={false}
            width={200}
            open={this.props.isOpenDrawer}
            onRequestChange={this.handleRequestChange}
          >
            <MenuItem onClick={this.handleClose}>
              <Link className="menu-list" to="/">
                Home
              </Link>
            </MenuItem>
            <MenuItem onClick={this.handleClose}>
              <Link className="menu-list" to="/about">
                About
              </Link>
            </MenuItem>
            <MenuItem onClick={this.handleClose}>
              <Link className="menu-list" to="/counter">
                Counter
              </Link>
            </MenuItem>
          </Drawer>

          <Route exact={true} path="/" component={Home} />
          <Route path="/about" component={About} />
          <Route path="/counter" component={ConnectedCounter} />
        </div>
      </Router>
    );
  }

  private handleLeftIconButtonTouchTap() {
    if (typeof this.props.actions !== 'undefined') {
      this.props.actions.toggleDrawer();
    }
  }

  private handleClose() {
    if (typeof this.props.actions !== 'undefined') {
      this.props.actions.toggleDrawer();
    }
  }

  private handleRequestChange(_: boolean) {
    if (typeof this.props.actions !== 'undefined') {
      this.props.actions.toggleDrawer();
    }
  }
}

function mapStateToProps(state: IRootState) {
  return state.app;
}

function mapDispatchToProps(dispatch: Dispatch<any>) {
  return { actions: bindActionCreators(appActionCreators, dispatch) };
}

export default connect(mapStateToProps, mapDispatchToProps)(App as any);
