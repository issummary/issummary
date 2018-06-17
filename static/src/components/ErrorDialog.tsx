import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';
import * as React from 'react';

interface IDialogAlertProps {
  open: boolean;
  error: string;
  onRequestClose: () => void;
}

export class ErrorDialog extends React.Component<IDialogAlertProps, any> {
  public render() {
    const actions = [
      <FlatButton
        key="OK"
        label="OK"
        primary={true}
        onClick={this.props.onRequestClose}
      />
    ];

    return (
      <Dialog
        title="Error occurred"
        actions={actions}
        modal={false}
        open={this.props.open}
        onRequestClose={this.props.onRequestClose}
      >
        {this.props.error}
      </Dialog>
    );
  }
}
